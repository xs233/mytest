#include <fstream>
#include "log.hpp"

Log::Log() : over_(false)
{

}

Log::~Log()
{
    over_ = true;
    cv_.notify_one();
    th_.join();
}

void Log::init()
{
    th_ = std::thread{&Log::write,this};
}

void Log::log(std::string& info)
{
    std::unique_lock<std::mutex> lck(mtx_);
    log_queue_.push(info);
    cv_.notify_one();
}

void Log::write()
{
    std::unique_lock<std::mutex> lck(mtx_);
    while (true)
    {
        cv_.wait(lck);
        if (over_)
            break;

        std::ofstream file("/home/spark/dolphin/bin/driver.log",std::ios::app | std::ios::binary);
        if (!file)
        {
            continue;
        }
        
        if (log_queue_.empty())
        {
            file << "there is no msg to write,why call me up?"  << '\n';
            file.close();
            continue;
        }

        do
        {
            file << log_queue_.front()  << '\n';
            log_queue_.pop();
        } while (!log_queue_.empty());
        file.close();
    }
}

