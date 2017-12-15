#ifndef __COMMAND_FACTORY_H__
#define __COMMAND_FACTORY_H__
class MonitorCommand;
class TaskCommand;

class CommandFactory
{
public:
    CommandFactory() {};
    virtual ~CommandFactory() {};
    virtual MonitorCommand* create_monitor_inst() = 0;
    virtual TaskCommand* create_task_inst() = 0;

};

#endif
