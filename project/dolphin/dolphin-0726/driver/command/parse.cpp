#include "parse.hpp"

Parse::Parse() : over_(false), logger_(LoggerPtr(Logger::getLogger("parse")))
{

}

Parse::~Parse()
{
    over_ = true;
    cv_.notify_one();
    th_.join();
    LOG4CXX_INFO(logger_,"parse over...");
}

void Parse::init(std::function<bool(int,std::string&)> fun_rule)
{
    th_ = std::thread{&Parse::parse,this};
    fun_rule_ = fun_rule;
}

void Parse::push(int fd,std::string& msg)
{
    std::unique_lock<std::mutex> lck(mtx_);
    socket_fd_ = fd;
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
            LOG4CXX_INFO(logger_,"command:" << str_command_);
            if (!fun_rule_)
            {
                LOG4CXX_ERROR(logger_,"fun_rule_ is empty!");
                break;
            }
            if (!fun_rule_(socket_fd_,str_command_))
                LOG4CXX_INFO(logger_,"unknown format");
        } while(!queue_empty);
        lck.lock();
    }
}

