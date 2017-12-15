#ifndef __LOG_H__
#define __LOG_H__

#include <string>
#include <queue>

class Log
{
public:
    static Log& instance()
    {
        static Log inst;
        return inst;
    };
    void write(std::string info);
    

private:
    Log();
    ~Log();
    std::queue<std::string> log_queue_;

};

#endif
