#include <iostream>
#define CRYPTOPP_ENABLE_NAMESPACE_WEAK 1
#include <cryptopp/md5.h>
#include <cryptopp/hex.h>
#include "task_command.hpp"

TaskCommand::TaskCommand() : Command("TaskCommand")
{

}

TaskCommand::~TaskCommand()
{

}

void TaskCommand::run_command(json& json_command)
{
   std::cout << "json:" << json_command << '\n';
}
