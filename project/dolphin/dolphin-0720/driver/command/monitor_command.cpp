#include <iostream>
#include "monitor_command.hpp"

MonitorCommand::MonitorCommand()
{

}

MonitorCommand::~MonitorCommand()
{

}

bool MonitorCommand::parse_command(std::string& str_command)
{
    str_command_ = str_command;
    return true;
}

bool MonitorCommand::run_command()
{
    std::cout << "run command:" << str_command_ << '\n';
    return true;
}

