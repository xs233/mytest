#ifndef __COMMAND_H__
#define __COMMAND_H__
#include <string>
#define CRYPTOPP_ENABLE_NAMESPACE_WEAK 1
#include <cryptopp/md5.h>
#include <cryptopp/hex.h>
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"
#include "../../json/json.hpp"
using namespace log4cxx;
using json = nlohmann::json;

class Command
{
public:
    Command(const char* log_name) : logger_(LoggerPtr(Logger::getLogger(log_name))) {};
    virtual ~Command() {};
    virtual bool run_command(json& json_command) = 0;
    virtual bool write_result(int fd) = 0;

protected:
    LoggerPtr logger_;
    std::string result_;
};

#endif
