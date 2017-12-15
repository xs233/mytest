#include <unistd.h>
#include <iostream>
#include "monitor_command.hpp"
#include "../config/config.hpp"

MonitorCommand::MonitorCommand()
{

}

MonitorCommand::~MonitorCommand()
{

}

void MonitorCommand::run_command(int fd,std::string& str_command)
{
    fd_ = fd;
    if (str_command.find("nodes") != std::string::npos)
    {
        std::string slave_info;
        Config::instance().get_slave(slave_info);
        if (slave_info.empty())
        {
            std::string error{"nodes is empty,it's unbelieveable"};
            write_fd(error);
        }
        else
        {
            write_fd(slave_info);
        }
    }
    else
    {
        std::string error{"unknown command"};
        write_fd(error);
    }
}

bool MonitorCommand::write_fd(std::string& str)
{
    if (write(fd_,str.c_str(),str.size()) == -1)
    {
        LOG4CXX_ERROR(logger_,"write monitor fd failed!");
        return false;
    }
    return true;
}
