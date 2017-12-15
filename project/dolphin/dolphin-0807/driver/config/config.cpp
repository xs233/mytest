#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fstream>
#include <ctime>
#include <chrono>
#include "config.hpp"

Config::Config() : logger_(LoggerPtr(Logger::getLogger("config"))), total_cores_(0)
{

}

Config::~Config()
{
    for (auto& kvp : slave_)
    {
        delete[] kvp.second;
    }
}

bool Config::start_executor()
{
    if (slave_.empty())
        return false;

    for (auto& kvp : slave_)
    {
        std::string shell_command = "ssh " + kvp.first + " \"cd  dolphin/bin; ./executor " + ip_master_ + " " + port_ + "\"";
	LOG4CXX_INFO(logger_,shell_command);
	int ret = system(shell_command.c_str());
        if (ret == -1)
	{
            LOG4CXX_ERROR(logger_,"ssh shell failed!");
            return false;
	}
    }
    return true;
}

bool Config::load_file()
{
    PropertyConfigurator::configure("driver.properties");
    std::ifstream file("node.conf",std::ios::in | std::ios::binary);
    if (!file)
    {
        LOG4CXX_ERROR(logger_,"open config file failed!");
	return false;
    }

    std::string ip;
    while (file >> ip)
    {
        size_t dot = ip.find(':');
        if (dot != std::string::npos)
	{
	    ip_master_ = ip.substr(0,dot);
	    port_ = ip.substr(dot + 1);
	}
	else
        {
            int* slave_info = new int[INFO_SIZE];
            memset(slave_info,'\0',INFO_SIZE * sizeof(int));
            slave_info[FD] = -1;
	    slave_.insert(std::pair<std::string,int*>(ip,slave_info));
        }
    }
    
    for (const auto& kvp : slave_)
        LOG4CXX_INFO(logger_,kvp.first << " status:waiting");
    file.close();
    return true;
}

int Config::get_port()
{
    return std::stoi(port_);
}

bool Config::set_slave(std::string& ip,int fd)
{
    std::unique_lock<std::mutex> lck{mtx_slave_};
    if (slave_.find(ip) == slave_.end())
        return false;
    slave_[ip][FD] = fd;
    if (fd == -1)
        total_cores_ -= slave_[ip][CORES];
    return true;
}

void Config::get_slave(std::string& slave_info)
{
    std::unique_lock<std::mutex> lck{mtx_slave_};
    for (const auto& kvp : slave_)
    {
        char begin_time[100] = {0};
        if (kvp.second[CT])
            std::strftime(begin_time,100,"%Y-%m-%d %X",std::localtime((std::time_t*)&kvp.second[CT]));
        else
            strcpy(begin_time,"no time info");
        std::string info = kvp.first + " fd:" + std::to_string(kvp.second[FD]) + "cores:" + std::to_string(kvp.second[CORES]) + "beginTime:" + begin_time + "\r\n";
        slave_info += info;
    }
}

bool Config::insert_client(int fd,std::string& ip)
{
    std::unique_lock<std::mutex> lck{mtx_client_};
    if (client_.find(fd) != client_.end())
    {
        LOG4CXX_ERROR(logger_,"insert_client failed!");
        close(fd);
        return false;
    }
    client_.insert(std::pair<int,std::string>(fd,ip));
    return true;
}

bool Config::set_slave(int fd,int cores)
{
    std::unique_lock<std::mutex> lck_client{mtx_client_};
    if (client_.find(fd) == client_.end())
    {
        LOG4CXX_ERROR(logger_,"find client failed");
        return false;
    }
    std::string ip = client_[fd];
    std::unique_lock<std::mutex> lck_slave{mtx_slave_};
    if (slave_.find(ip) == slave_.end())
    {
        LOG4CXX_ERROR(logger_,"find slave failed");
        return false;
    }
    if (slave_[ip][FD] != -1)
        return false;
    slave_[ip][FD] = fd;
    slave_[ip][CORES] = cores;
    total_cores_ += cores;
    using namespace std::chrono;
    slave_[ip][CT] = duration_cast< seconds >(system_clock::now().time_since_epoch()).count();
    return true;
}

bool Config::erase_client(int fd)
{
    std::unique_lock<std::mutex> lck{mtx_client_};
    if (client_.find(fd) == client_.end())
    {
        LOG4CXX_ERROR(logger_,"erase_client failed!");
        return false;
    }
    if (set_slave(client_[fd],-1))
        LOG4CXX_ERROR(logger_,"one node disconnected,please check it!");
    client_.erase(fd);
    return true;
}

void Config::get_client(std::string& client_info)
{
    std::unique_lock<std::mutex> lck{mtx_client_};
    for (const auto& kvp : client_)
    {
        std::string info = std::to_string(kvp.first) + ':' + kvp.second + "\r\n";
        client_info += info;
    }
}

int Config::get_total_cores()
{
    return total_cores_;
}

bool Config::assign_work(std::string& md5_value,int total_num,std::map<std::string,json>& col_info)
{
    int remaining_cores = total_cores_;
    if (!remaining_cores) {
        LOG4CXX_ERROR(logger_,"there is no vaild slave!");
        return false;
    }

    if (total_num < remaining_cores) {
        LOG4CXX_ERROR(logger_,"the total_num is " << total_num << ",and the total_cores is " << total_cores_);
        return false;
    }
    int per_core_task = (total_num - total_num % remaining_cores) / remaining_cores;
    std::unique_lock<std::mutex> lck{mtx_slave_};
    for (const auto& kvp : slave_) {
        int* slave_info = kvp.second;
        if (slave_info[FD] != -1) {
            if (slave_info[CORES] == remaining_cores) {
                if (!send_task_to_slave(md5_value,slave_info[FD],col_info,true,slave_info[CORES],per_core_task))
                    return false;
            } else {
                if (!send_task_to_slave(md5_value,slave_info[FD],col_info,false,slave_info[CORES],per_core_task))
                    return false;
            }
            remaining_cores -= slave_info[CORES];
        }       
    }
    return true;
}

bool Config::send_task_to_slave(std::string& md5_value,int fd,std::map<std::string,json>& col_info,bool last_one,int cores,int per_core_task)
{
    json slave_task;
    for (int index = 0; index < cores; ++index) {
        json json_per_core;
        bool last_core = false;
        if (index == (cores - 1))
            last_core = true;
        for (auto it = col_info.begin(); it != col_info.end(); ) {
            if (last_one && last_core) {
                json_per_core.push_back(it->second);
                ++it;
                continue;        
            }

            std::cout << "111" << '\n';
            json& time = it->second;
            std::cout << time << '\n';
            std::cout << "222" << '\n';
            long end_time = time["endTime"];
            std::cout << "333" << '\n';
            long begin_time = time["beginTime"];
            int col_num = end_time - begin_time + 1;
            std::cout << "444" << '\n';
            if (col_num > per_core_task) {
                // end_time = time["endTime"];
                time["endTime"] = begin_time + per_core_task - 1;
                json_per_core.push_back(time);
                // time["beginTime"] = (long)time["endTime"] + 1;
                time["beginTime"] = begin_time + per_core_task;
                time["endTime"] = end_time;
                ++it;
            } else {
                json_per_core.push_back(time);
                col_info.erase(it++);
            }
            per_core_task -= col_num;
            if (per_core_task <= 0)
                break;
        }
        slave_task.push_back(json_per_core);
    }

    std::string send_msg = slave_task.dump();
    int ret = write(fd,send_msg.c_str(),send_msg.length());
    if (ret == -1) {
        LOG4CXX_ERROR(logger_,"send task to slave failed!");
        return false;    
    }
    return true;
}
