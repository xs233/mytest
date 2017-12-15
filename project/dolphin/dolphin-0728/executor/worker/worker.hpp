#ifndef __WORKER_H__
#define __WORKER_H__
#include "../../json/json.hpp"
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"

using namespace log4cxx;
using json = nlohmann::json;

class Worker
{
public:
    Worker(const char* log_name) : logger_(LoggerPtr(Logger::getLogger(log_name))) {};
    virtual ~Worker() {};
    virtual void run(int fd,json& task_json) {};
protected:
    LoggerPtr logger_;

};

#endif
