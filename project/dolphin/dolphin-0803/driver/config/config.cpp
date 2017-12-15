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

bool Config::assign_work(int total_num,std::map<std::string,json>& col_info)
{
    int remaining_cores = total_cores_;
    if (!remaining_cores) {
        return false;
    }

    int per_core_task = (total_num - total_num % remaining_cores) / remaining_cores;
    std::unique_lock<std::mutex> lck{mtx_slave_};
    for (const auto& kvp : slave_) {

    }
}

