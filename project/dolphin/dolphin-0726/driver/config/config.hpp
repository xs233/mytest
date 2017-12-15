#ifndef __CONFIG_HPP__
#define __CONFIG_HPP__

#include <string>
#include <map>
#include <mutex>
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"
using namespace log4cxx;

class Config
{
public:
    static Config& instance()
    {
        static Config inst;
        return inst;
    };
    bool load_file();
    bool start_executor();
    int get_port();
    void get_slave(std::string& slave_info);
    bool insert_client(int fd,std::string& ip);
    bool set_client(int fd);
    bool erase_client(int fd);
    void get_client(std::string& client_info);
    Config(const Config& ) = delete;
    Config& operator=(const Config& ) = delete;

private:
    Config();
    ~Config();
    bool set_slave(std::string& ip,int fd);
    std::map<std::string,int> slave_;
    std::map<int,std::string> client_;
    std::string ip_master_;
    std::string port_;
    LoggerPtr logger_;
    std::mutex mtx_slave_;
    std::mutex mtx_client_;
};

#endif
