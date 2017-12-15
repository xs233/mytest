#ifndef __PARSE_H__
#define __PARSE_H__
#include <string>
#include <queue>
#include <map>
#include <mutex>
#include <condition_variable>
#include <thread>
#include "../../json/json.hpp"
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"

class CommandFactoryInst;
using json = nlohmann::json;
using namespace log4cxx;

class Parse
{
public:
    Parse(const Parse& ) = delete;
    Parse& operator=(const Parse& ) = delete;
    static Parse& instance()
    {
        static Parse s_parse;
        return s_parse;
    };
    void init();
    void push(int fd,std::string& msg);

private:
    Parse();
    ~Parse();
    void parse();
    //void flag_parse();
    void md5(std::string& md5_str,std::string& md5_value);
    bool insert_task(int fd,std::string& str,json& json_task);
    enum
    {
        FD = 0,
        STATUS = 1,
        GT = 2,    //GetTime for data
        ST = 3,    //StartTime
        ET = 4     //EndTime
    };
    enum
    {
        WAITTING = 0,
        RUNNING = 1,
        END = 2,
        FD_ERROR = 3
    };
    static const int INFO_SIZE = 10;
    bool over_;
    std::thread th_;
    std::mutex mtx_;
    std::mutex mtx_fd_task_;
    std::mutex mtx_task_info_;
    std::condition_variable cv_;
    std::queue<json> command_queue_;
    std::map<int,std::string> fd_task_;
    std::map<std::string,int*> task_info_;
    LoggerPtr logger_;
};

#endif
