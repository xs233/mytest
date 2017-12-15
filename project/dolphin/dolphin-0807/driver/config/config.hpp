#ifndef __CONFIG_HPP__
#define __CONFIG_HPP__

#include <string>
#include <map>
#include <mutex>
#include "../../json/json.hpp"
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"
using namespace log4cxx;
using json = nlohmann::json;

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
    bool set_slave(int fd,int cores);
    bool erase_client(int fd);
    void get_client(std::string& client_info);
    int get_total_cores();
    bool assign_work(std::string& md5_value,int total_num,std::map<std::string,json>& col_info);
    Config(const Config& ) = delete;
    Config& operator=(const Config& ) = delete;

private:
    Config();
    ~Config();
    bool set_slave(std::string& ip,int fd);
    bool send_task_to_slave(std::string& md5_value,int fd,std::map<std::string,json>& col_info,bool last_one,int cores,int per_core_task);
    enum
    {
        FD = 0,
        CORES = 1,
        CT = 2  //ConnectTime
    };
    static const int INFO_SIZE = 5;
    std::map<std::string,int*> slave_;
    std::map<int,std::string> client_;
    std::string ip_master_;
    std::string port_;
    LoggerPtr logger_;
    std::mutex mtx_slave_;
    std::mutex mtx_client_;
    int total_cores_;
};

#endif
