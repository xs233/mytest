#ifndef __MONITOR_COMMAND_H__
#define __MONITOR_COMMAND_H__
#include "command.hpp"

class MonitorCommand : public Command
{
public:
    MonitorCommand();
    ~MonitorCommand();
    virtual bool parse_command(std::string& str_command) override;
    virtual bool run_command() override;

private:
    std::string str_command_;
};

#endif
