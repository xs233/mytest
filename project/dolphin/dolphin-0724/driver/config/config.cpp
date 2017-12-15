#include <stdlib.h>
#include <fstream>
#include "config.hpp"

Config::Config() : logger_(LoggerPtr(Logger::getLogger("config")))
{

}

Config::~Config()
{
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
	    slave_.insert(std::pair<std::string,int>(ip,-1));
    }
    
    for (const auto& kvp : slave_)
        LOG4CXX_INFO(logger_,kvp.first << " status:unknown");
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
    {
        LOG4CXX_ERROR(logger_,"set_slave failed!");
        return false;
    }
    slave_[ip] = fd;
    return true;
}

void Config::get_slave(std::string& slave_info)
{
    std::unique_lock<std::mutex> lck{mtx_slave_};
    for (const auto& kvp : slave_)
    {
        std::string info = kvp.first + " status:" + std::to_string(kvp.second) + "\r\n";
        slave_info += info;
    }
}

bool Config::insert_client(int fd,std::string& ip)
{
    std::unique_lock<std::mutex> lck{mtx_client_};
    if (client_.find(fd) != client_.end())
    {
        LOG4CXX_ERROR(logger_,"insert_client failed!");
        return false;
    }
    client_.insert(std::pair<int,std::string>(fd,ip));
    return true;
}

bool Config::set_client(int fd)
{
    std::unique_lock<std::mutex> lck{mtx_client_};
    if (client_.find(fd) == client_.end())
    {
        LOG4CXX_ERROR(logger_,"set_client failed");
        return false;
    }
    if (!set_slave(client_[fd],fd))
        return false;
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

