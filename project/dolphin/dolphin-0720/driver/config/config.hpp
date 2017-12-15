#ifndef __CONFIG_HPP__
#define __CONFIG_HPP__

#include <string>
#include <vector>
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
    Config(const Config& ) = delete;
    Config& operator=(const Config& ) = delete;

private:
    Config();
    ~Config();
    std::vector<std::string> ip_slave_;
    std::string ip_master_;
    std::string port_;
    LoggerPtr logger_;

};

#endif
