#ifndef __DRIVER_HPP__
#define __DRIVER_HPP__
#include<sys/epoll.h>
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"

using namespace log4cxx;

class Driver
{
public:
    Driver();
    ~Driver();
    bool start_driver();

private:
    static const int CONN_MAX = 10;
    static const int BUF_SIZE = 10240;
    epoll_event evs_[CONN_MAX];
    int socket_fd_;
    int epfd_;
    LoggerPtr logger_;
    char* recv_buf_;
    bool init_tcp_env();
    bool epoll();
    int run_command(int cnt);
};

#endif
