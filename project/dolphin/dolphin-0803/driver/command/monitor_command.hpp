#ifndef __MONITOR_COMMAND_H__
#define __MONITOR_COMMAND_H__
#include "command.hpp"

class MonitorCommand : public Command
{
public:
    MonitorCommand();
    ~MonitorCommand();
    virtual bool run_command(json& json_command) override;
    virtual bool write_result(int fd) override;

};

#endif
