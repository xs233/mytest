#ifndef __LOG_H__
#define __LOG_H__

#include <string>
#include <queue>
#include <mutex>
#include <condition_variable>
#include <thread>

class Log
{
public:
    static Log& instance()
    {
        static Log inst;
        return inst;
    };
    void log(std::string& info);
    void init();
    Log(const Log& ) = delete;
    Log& operator=(const Log& ) = delete;

private:
    Log();
    ~Log();
    void write();

    bool over_;
    std::thread th_;
    std::mutex mtx_;
    std::condition_variable cv_;
    std::queue<std::string> log_queue_;

};

#endif
