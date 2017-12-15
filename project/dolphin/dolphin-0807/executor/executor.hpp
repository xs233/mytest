#ifndef __EXECUTOR_H__
#define __EXECUTOR_H__
#include <string>
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h" 

using namespace log4cxx;

class Executor
{
public:
    Executor(char* master_ip,int port);
    ~Executor();
    bool init();
    bool process();
    

private:
    bool recv_msg();
    bool send_msg(std::string& msg);
    static const int BUF_SIZE = 10240;
    std::string master_ip_;
    int port_;
    int socket_fd_;
    LoggerPtr logger_;
    char* recv_buf_;
};

#endif
