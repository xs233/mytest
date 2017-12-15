#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include "executor.hpp"
#include "worker/parse.hpp"
#include "worker/worker_factory_inst.hpp"
#include "worker/mongodb_worker.hpp"

bool flag_parse(int fd,std::string& str_command)
{
    json json_command;
    try
    {
        json_command = json::parse(str_command);
        if (json_command.find("task") == json_command.end())
            throw 1;
 
        json flag = json_command["task"];
        if (!flag.is_string())
            throw 1;
        else if (flag == "mongodb")
        {
            WorkerFactoryInst inst;
            MongodbWorker* tmp = inst.create_mongodb_inst();
            tmp->run(fd,json_command);
            delete tmp;
        }
        else
            throw 1;
    }
    catch (...)
    {
        return false;
    }
    return true;
}

Executor::Executor(char* master_ip,int port) : master_ip_(master_ip), port_(port), logger_(LoggerPtr(Logger::getLogger("executor"))), recv_buf_(new char[BUF_SIZE])
{

}

Executor::~Executor()
{
    if (recv_buf_)
    {
        delete[] recv_buf_;
        recv_buf_ = nullptr;
    }
}

bool Executor::init()
{
    PropertyConfigurator::configure("/home/spark/dolphin/bin/executor.properties");
    Parse::instance().init(flag_parse);
    LOG4CXX_INFO(logger_,"init start");
    socket_fd_ = socket(PF_INET,SOCK_STREAM,0);
    if (socket_fd_ == -1)
    {
        LOG4CXX_ERROR(logger_, "socket failed!");
        return false;
    }

    sockaddr_in addr = {0};
    socklen_t addrlen = sizeof(addr);
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = inet_addr(master_ip_.c_str());
    addr.sin_port = htons(port_);

    try
    {
        int ret = connect(socket_fd_,(struct sockaddr*)&addr,addrlen);
        if (ret == -1)
            throw "connect failed!";
        if (!recv_msg())
            throw "init first recv msg failed!";
        if (strcmp(recv_buf_,"welcome"))
            throw "the msg is not the welcome!";
        std::string msg = "dolphin slave," + std::to_string(sysconf(_SC_NPROCESSORS_ONLN));
        LOG4CXX_INFO(logger_,"msg:" << msg);
        if (!send_msg(msg))
            throw "init send msg failed!";
        if (!recv_msg())
            throw "init second recv msg failed!";
        if (strcmp(recv_buf_,"connect success"))
            throw "the msg is not the connect success";
    }
    catch (const char* s)
    {
        LOG4CXX_ERROR(logger_,s);
        close(socket_fd_);
        return false;
    }
    LOG4CXX_INFO(logger_,"executor init success...");
    return true;
}

bool Executor::process()
{
    LOG4CXX_INFO(logger_,"wait for task...");
    while (true)
    {
        if (!recv_msg())
        {
            LOG4CXX_ERROR(logger_,"recv task failed!");
        }
        else
        {
            if (!strcmp(recv_buf_,"executor exit"))
            {
                LOG4CXX_INFO(logger_,"executor exit");
                close(socket_fd_);
                return true;
            }
            std::string msg{recv_buf_};
            LOG4CXX_INFO(logger_,msg);
            Parse::instance().push(socket_fd_,msg);
            LOG4CXX_INFO(logger_,"push msg");
        }
    }
    return true;
}

bool Executor::recv_msg()
{
    memset(recv_buf_,'\0',BUF_SIZE);
    int ret = recv(socket_fd_,recv_buf_,BUF_SIZE,0);
    if (ret == -1)
        return false;
    return true;
}

bool Executor::send_msg(std::string& msg)
{
    int ret = send(socket_fd_,msg.c_str(),msg.size(),0);
    if (ret == -1)
        return false;
    return true;
}
