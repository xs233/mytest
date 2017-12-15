#include <iostream>
#include "task_command.hpp"

TaskCommand::TaskCommand()
{

}

TaskCommand::~TaskCommand()
{

}

bool TaskCommand::parse_command(std::string& str_command)
{
    str_command_ = str_command;
    return true;
}

bool TaskCommand::run_command()
{
    std::cout << "run command:" << str_command_ << '\n';
    return true;
}

