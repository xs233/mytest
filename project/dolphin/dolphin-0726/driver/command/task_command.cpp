#include <iostream>
#include "task_command.hpp"

TaskCommand::TaskCommand()
{

}

TaskCommand::~TaskCommand()
{

}

void TaskCommand::run_command(int fd,json& json_command)
{
    std::cout << fd << json_command << '\n';
}
