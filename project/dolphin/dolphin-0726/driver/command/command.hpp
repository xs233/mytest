#ifndef __COMMAND_H__
#define __COMMAND_H__
#include <string>
#include "../../json/json.hpp"
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"
using namespace log4cxx;
using json = nlohmann::json;

class Command
{
public:
    Command() : logger_(LoggerPtr(Logger::getLogger("command"))) {};
    virtual ~Command() {};
    virtual void run_command(int fd,std::string& str_command) {};
    virtual void run_command(int fd,json& json_command) {};

protected:
    LoggerPtr logger_;
};

#endif
