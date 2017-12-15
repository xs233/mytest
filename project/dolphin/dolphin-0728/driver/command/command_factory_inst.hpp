#ifndef __COMMAND_FACTORY_INST_H__
#define __COMMAND_FACTORY_INST_H__

#include "command_factory.hpp"

class CommandFactoryInst : public CommandFactory
{
public:
    CommandFactoryInst();
    ~CommandFactoryInst();
    virtual MonitorCommand* create_monitor_inst() override;
    virtual TaskCommand* create_task_inst() override;

};

#endif
