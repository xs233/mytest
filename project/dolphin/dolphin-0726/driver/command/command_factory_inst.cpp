#include "monitor_command.hpp"
#include "task_command.hpp"
#include "command_factory_inst.hpp"

CommandFactoryInst::CommandFactoryInst()
{

}

CommandFactoryInst::~CommandFactoryInst()
{

}

MonitorCommand* CommandFactoryInst::create_monitor_inst()
{
    return new MonitorCommand{};
}

TaskCommand* CommandFactoryInst::create_task_inst()
{
    return new TaskCommand{};
}

