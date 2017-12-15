#ifndef __TASK_COMMAND_H__
#define __TASK_COMMAND_H__
#include "command.hpp"

class TaskCommand : public Command
{
public:
    TaskCommand();
    ~TaskCommand();
    virtual void run_command(int fd,json& json_command) override;

};

#endif
