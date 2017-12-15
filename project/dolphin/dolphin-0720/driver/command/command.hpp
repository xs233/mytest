#ifndef __COMMAND_H__
#define __COMMAND_H__
#include <string>

class Command
{
public:
    Command() {};
    virtual ~Command() {};
    virtual bool parse_command(std::string& str_command) = 0;
    virtual bool run_command() = 0;

};

#endif
