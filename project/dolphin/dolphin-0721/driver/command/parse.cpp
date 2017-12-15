#include <unistd.h>
#include "parse.hpp"
#include "command_factory_inst.hpp"
#include "monitor_command.hpp"
#include "task_command.hpp"
#include "../config/config.hpp"

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
            LOG4CXX_INFO(logger_,"command:" << str_command_);
            flag_parse();
        } while(!queue_empty);
        lck.lock();
    }
}

void Parse::flag_parse()
{
    size_t pos = str_command_.find("||");
    if (pos == std::string::npos || pos == 0)
    {
        LOG4CXX_ERROR(logger_,"this is not the correct format,why it happened!");
    }

    int fd = std::stoi(str_command_.substr(0,pos));
    std::string info{str_command_.substr(pos+2)};
    if (!info.compare(0,7,"dolphin"))
    {
        MonitorCommand* tmp = inst_->create_monitor_inst();
        tmp->run_command(fd,info);
        delete tmp;
    }
    else
    {
        
        json json_command;
        try
        {
            json_command = json::parse(info);
            if (json_command.find("task") == json_command.end())
            {
                std::cout << "not find" << '\n';
                throw 1;
            }

            json flag = json_command["task"];
            if (!flag.is_string())
                throw 1;
            else if (flag == "mongodb")
            {
                TaskCommand* tmp = inst_->create_task_inst();
                tmp->run_command(fd,json_command);
                delete tmp;
            }
            else
                throw 1;
        }
        catch (...)
        {
            std::string er{"unknown format"};
            write(fd,er.c_str(),er.size());
        }
    }
}

