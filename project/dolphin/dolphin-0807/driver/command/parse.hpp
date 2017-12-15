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

class CommandFactory;
using json = nlohmann::json;
using namespace log4cxx;

enum CommandObject {ERROR,UNKNOWN,MONITOR,MONGODB};
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
    bool push(int fd,char* msg,int len);
    void remove_task(int fd,bool fd_close);

private:
    Parse();
    ~Parse();
    void worker();
    void md5(char* md5_str,int len,std::string& md5_value);
    bool get_msg(std::string& md5_value,std::string& msg,json& json_command);
    int get_fd(std::string& md5_value);
    CommandObject parse_command(std::string& md5_value,json& json_command);
    void do_work(std::string& md5_value);
    enum
    {
        FD = 0,
        MSG = 1,
        GT = 2,    //GetTime for data
        ST = 3,    //StartTime
        ET = 4     //EndTime
    };
    static const int INFO_SIZE = 10;
    bool over_;
    std::thread th_;
    std::mutex mtx_;
    std::mutex mtx_fd_task_;
    std::mutex mtx_task_info_;
    std::condition_variable cv_;
    std::queue<std::string> command_queue_;
    std::map<int,std::string> fd_task_;
    std::map<std::string,long*> task_info_;
    LoggerPtr logger_;
    CommandFactory* inst_;
};

#endif
