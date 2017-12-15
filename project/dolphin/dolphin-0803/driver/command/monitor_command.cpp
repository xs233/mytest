#include <unistd.h>
#include <iostream>
#include "monitor_command.hpp"
#include "../config/config.hpp"

MonitorCommand::MonitorCommand() : Command("MonitorCommand")
{

}

MonitorCommand::~MonitorCommand()
{

}

bool MonitorCommand::run_command(json& json_command)
{
    std::string str_command = json_command["MSG"];
    if (str_command.find("nodes") != std::string::npos)
    {
        std::string slave_info;
        Config::instance().get_slave(slave_info);
        if (slave_info.empty())
        {
            result_ = "nodes is empty,it's unbelieveable";
        }
        else
        {
            result_ = slave_info;
        }
    }
    else if (str_command.find("slave,") != std::string::npos)
    {
        int cores = std::stoi(str_command.substr(str_command.find(',') + 1,1));
        int fd = json_command["FD"];
        if (Config::instance().set_slave(fd,cores))
        {
            result_ = "connect success";
        }
        else
        {
            result_ = "unknown error";
        }
    }
    else if (str_command.find("clients") != std::string::npos)
    {
        std::string client_info;
        Config::instance().get_client(client_info);
        result_ = client_info;
    }
    else
    {
        result_ = "unknown command";
    }
    return true;
}

bool MonitorCommand::write_result(int fd)
{
    if (write(fd,result_.c_str(),result_.size()) == -1)
    {
        LOG4CXX_ERROR(logger_,"write monitor fd failed!");
        return false;
    }
    return true;
}

