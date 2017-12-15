#ifndef __TASK_COMMAND_H__
#define __TASK_COMMAND_H__
#include "command.hpp"

class TaskCommand : public Command
{
public:
    TaskCommand();
    ~TaskCommand();
    virtual bool parse_command(std::string& str_command) override;
    virtual bool run_command() override;

private:
    std::string str_command_;
};

#endif
