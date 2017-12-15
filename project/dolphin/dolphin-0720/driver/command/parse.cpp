#include <string.h>
#include <unistd.h>
#include "parse.hpp"
#include "command_factory_inst.hpp"
#include "monitor_command.hpp"
#include "task_command.hpp"

Parse::Parse() : inst_(new CommandFactoryInst{}), over_(false), logger_(LoggerPtr(Logger::getLogger("parse")))
{

}

Parse::~Parse()
{
    over_ = true;
    cv_.notify_one();
    th_.join();

    if (inst_)
    {
        delete inst_;
        inst_ = nullptr;
    }
    LOG4CXX_INFO(logger_,"parse over...");
}

void Parse::init()
{
    th_ = std::thread{&Parse::parse,this};
}

void Parse::push(std::string& msg)
{
    std::unique_lock<std::mutex> lck(mtx_);
    command_queue_.push(msg);
    lck.unlock();
    cv_.notify_one();
}

void Parse::parse()
{
    std::unique_lock<std::mutex> lck(mtx_);
    while (true)
    {
        cv_.wait(lck);
        if (over_)
            break;
        
        if (command_queue_.empty())
        {
            LOG4CXX_ERROR(logger_,"the command_queue is empty,why call me up?");
            continue;
        }

        bool queue_empty;
        bool first_loop = true;
        //lck.unlock();
        do
        {
            if (!first_loop)
                lck.lock();
            str_command_ = command_queue_.front();
            command_queue_.pop();
            queue_empty = command_queue_.empty();
            lck.unlock();
            flag_parse();
        } while(!queue_empty);
        lck.lock();
    }
}

void Parse::flag_parse()
{
    size_t pos = str_command_.find("||");
    if (pos == std::string::npos)
    {
        LOG4CXX_INFO(logger_,"this is not the correct format,why it happened!");
    }
    else if (str_command_.compare(pos+2,7,"dolphin"))
    {
        write(std::stoi(str_command_.substr(0,pos)),"unvaild command",strlen("unvaild command"));
    }
    LOG4CXX_INFO(logger_,"command:" << str_command_);
}

