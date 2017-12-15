#include <unistd.h>
#include <chrono>
#include "parse.hpp"
#define CRYPTOPP_ENABLE_NAMESPACE_WEAK 1
#include "cryptopp/md5.h"
#include "cryptopp/hex.h"
#include "command_factory_inst.hpp"
#include "monitor_command.hpp"
#include "task_command.hpp"


Parse::Parse() : over_(false), logger_(LoggerPtr(Logger::getLogger("parse")))
{

}

Parse::~Parse()
{
    over_ = true;
    cv_.notify_one();
    th_.join();
    for (auto& kvp : task_info_)
        delete[] kvp.second;
    LOG4CXX_INFO(logger_,"parse over...");
}

void Parse::init()
{
    th_ = std::thread{&Parse::parse,this};
}

void Parse::push(int fd,std::string& str_command)
{
    if (!str_command.compare(0,7,"dolphin"))
    {
        CommandFactoryInst inst;
        MonitorCommand* tmp = inst.create_monitor_inst();
        tmp->run_command(fd,str_command);
        delete tmp;
    }
    else
    {
        json json_command;
        try
        {
            json_command = json::parse(str_command);
            if (json_command.find("task") == json_command.end())
                throw 1;

            json flag = json_command["task"];
            if (!flag.is_string())
                throw 1;
            else if (flag == "mongodb")
            {
                if (insert_task(fd,str_command,json_command))
                    cv_.notify_one();
            }
            else
                throw 1;
        }
        catch (...)
        {
            std::string er{"unknown format"};
            write(fd,er.c_str(),er.size());
            return;
        }
    }
}

void Parse::md5(std::string& md5_str,std::string& md5_value)
{
    byte digest[ CryptoPP::Weak::MD5::DIGESTSIZE ];
    CryptoPP::Weak::MD5 hash;
    hash.CalculateDigest( digest, (const byte*)md5_str.c_str(), md5_str.length() );
    CryptoPP::HexEncoder encoder;
    encoder.Attach( new CryptoPP::StringSink( md5_value ) );
    encoder.Put( digest, sizeof(digest) );
    encoder.MessageEnd();
}

bool Parse::insert_task(int fd,std::string& str,json& json_task)
{
    std::string md5_value;
    md5(str,md5_value);
    std::unique_lock<std::mutex> lck_fd_task(mtx_fd_task_);
    if (fd_task_.find(fd) != fd_task_.end())
    {
        LOG4CXX_ERROR(logger_,"fd_task_.find failed!");
        return false;
    }
    fd_task_.insert(std::pair<int,std::string>(fd,md5_value));
    int* info = new int[INFO_SIZE];
    memset(info,'\0',INFO_SIZE * sizeof(int));
    info[FD] = fd;
    info[STATUS] = WAITTING;
    using namespace std::chrono;
    info[GT] = duration_cast< milliseconds >(system_clock::now().time_since_epoch()).count();
    std::unique_lock<std::mutex> lck_task_info(mtx_task_info_);
    if (task_info_.find(md5_value) != task_info_.end())
    {
        LOG4CXX_ERROR(logger_,"fd_info_.find failed!");
        delete[] info;
        return false;
    }
    task_info_.insert(std::pair<std::string,int*>(md5_value,info));
    std::unique_lock<std::mutex> lck_queue(mtx_);
    json_task["MD5"] = md5_value;
    command_queue_.push(json_task);
    return true;
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
            {
                first_loop = false;
                lck.lock();
            }
            json json_task = command_queue_.front();
            command_queue_.pop();
            queue_empty = command_queue_.empty();
            lck.unlock();
            LOG4CXX_INFO(logger_,"command:" << json_task);
            CommandFactoryInst inst;
            TaskCommand* tmp = inst.create_task_inst();
            tmp->run_command(json_task);
            delete tmp;
        } while(!queue_empty);
        lck.lock();
    }
}

