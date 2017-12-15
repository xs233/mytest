#ifndef __MONITOR_COMMAND_H__
#define __MONITOR_COMMAND_H__
#include "command.hpp"

class MonitorCommand : public Command
{
public:
    MonitorCommand();
    ~MonitorCommand();
    virtual void run_command(int fd,std::string& str_command) override;

private:
    bool write_fd(std::string& str);
    int fd_;
};

#endif
