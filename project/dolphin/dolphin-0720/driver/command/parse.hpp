#ifndef __PARSE_H__
#define __PARSE_H__
#include <string>
#include <queue>
#include <mutex>
#include <condition_variable>
#include <thread>
#include "log4cxx/logger.h"
#include "log4cxx/propertyconfigurator.h"

class CommandFactoryInst;
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
    void push(std::string& msg);

private:
    Parse();
    ~Parse();
    void parse();
    void flag_parse();
    CommandFactoryInst* inst_;
    bool over_;
    std::thread th_;
    std::mutex mtx_;
    std::condition_variable cv_;
    std::queue<std::string> command_queue_;
    LoggerPtr logger_;
    std::string str_command_;
};

#endif
