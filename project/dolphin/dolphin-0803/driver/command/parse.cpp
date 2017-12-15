#include <unistd.h>
#include <chrono>
#include "parse.hpp"
#include "command_factory_inst.hpp"
#include "monitor_command.hpp"
#include "task_command.hpp"


Parse::Parse() : over_(false), logger_(LoggerPtr(Logger::getLogger("parse"))), inst_(new CommandFactoryInst())
{

}

Parse::~Parse()
{
    over_ = true;
    cv_.notify_one();
    th_.join();
    std::unique_lock<std::mutex> lck(mtx_task_info_);
    for (auto& kvp : task_info_)
    {
        delete[] (char*)kvp.second[MSG];
        delete[] kvp.second;
    }
    if (inst_)
    {
        delete inst_;
        inst_ = nullptr;
    }
    LOG4CXX_INFO(logger_,"parse over...");
}

void Parse::init()
{
    th_ = std::thread{&Parse::worker,this};
}

bool Parse::push(int fd,char* msg,int len)
{
    std::string md5_value;
    md5(msg,len,md5_value);
    std::unique_lock<std::mutex> lck_fd_task(mtx_fd_task_);
    if (fd_task_.find(fd) != fd_task_.end())
    {
        LOG4CXX_ERROR(logger_,"fd_task_.find failed!");
        return false;
    }
    fd_task_.insert(std::pair<int,std::string>(fd,md5_value));
    long* info = new long[INFO_SIZE];
    memset(info,'\0',INFO_SIZE * sizeof(long));
    info[FD] = fd;
    char* pmsg = new char[len + 1];
    memset(pmsg,'\0',len+1);
    std::strcpy(pmsg,msg);
    std::cout << pmsg << '\n';
    info[MSG] = (long)pmsg;
    using namespace std::chrono;
    info[GT] = duration_cast< milliseconds >(system_clock::now().time_since_epoch()).count();
    std::unique_lock<std::mutex> lck_task_info(mtx_task_info_);
    if (task_info_.find(md5_value) != task_info_.end())
    {
        LOG4CXX_ERROR(logger_,"fd_info_.find failed!");
        delete[] info;
        return false;
    }
    task_info_.insert(std::pair<std::string,long*>(md5_value,info));
    std::unique_lock<std::mutex> lck_queue(mtx_);
    command_queue_.push(md5_value);
    cv_.notify_one();
    return true;
}

void Parse::remove_task(int fd,bool fd_close)
{
    std::unique_lock<std::mutex> lck_fd_task(mtx_fd_task_);
    auto it_fd_task = fd_task_.find(fd);
    if (it_fd_task == fd_task_.end())
        return;
    std::unique_lock<std::mutex> lck_task_info(mtx_task_info_);
    auto it_task_info = task_info_.find(fd_task_[fd]);
    if (it_task_info == task_info_.end())
        return;
    long* info = task_info_[fd_task_[fd]];
    using namespace std::chrono;
    info[ET] = duration_cast< milliseconds >(system_clock::now().time_since_epoch()).count();
    if (fd_close)
    {
        LOG4CXX_WARN(logger_,"FD_ERROR::" << "msg:" << (char*)info[MSG] << ",GetTime:" << info[GT] << ",BeginTime:" << info[ST] << ",EndTime:" << info[ET]);
    }
    else
    {
        LOG4CXX_INFO(logger_,"NORMAL::" << "msg:" << (char*)info[MSG] << ",GetTime:" << info[GT] << ",BeginTime:" << info[ST] << ",EndTime:" << info[ET]);
    }
    delete[] (char*)info[MSG];
    task_info_.erase(it_task_info);
    fd_task_.erase(it_fd_task);
}

void Parse::md5(char* md5_str,int len,std::string& md5_value)
{
    byte digest[ CryptoPP::Weak::MD5::DIGESTSIZE ];
    CryptoPP::Weak::MD5 hash;
    hash.CalculateDigest( digest, (const byte*)md5_str, len);
    CryptoPP::HexEncoder encoder;
    encoder.Attach( new CryptoPP::StringSink( md5_value ) );
    encoder.Put( digest, sizeof(digest) );
    encoder.MessageEnd();
}

void Parse::worker()
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
            std::string md5_value = command_queue_.front();
            command_queue_.pop();
            lck.unlock();
            do_work(md5_value);
            lck.lock();
            queue_empty = command_queue_.empty();
            lck.unlock();
        } while(!queue_empty);
        lck.lock();
    }
}

bool Parse::get_msg(std::string& md5_value,std::string& msg,json& json_msg)
{
    std::unique_lock<std::mutex> lck(mtx_task_info_);
    if (task_info_.find(md5_value) == task_info_.end())
    {
        LOG4CXX_ERROR(logger_,"the task has been deleted!");
        return false;
    }
    long* info = task_info_[md5_value];
    json_msg["FD"] = info[FD];
    msg = (char*)info[MSG];
    json_msg["MSG"] = (char*)info[MSG];
    using namespace std::chrono;
    info[ST] = duration_cast< milliseconds >(system_clock::now().time_since_epoch()).count();
    return true;
}

int Parse::get_fd(std::string& md5_value)
{
    std::unique_lock<std::mutex> lck(mtx_task_info_);
    if (task_info_.find(md5_value) == task_info_.end())
    {
        LOG4CXX_ERROR(logger_,"the task has been deleted!");
        return -1;
    }
    return task_info_[md5_value][FD];
}

CommandObject Parse::parse_command(std::string& md5_value,json& json_task)
{
    std::string info;
    bool ret = get_msg(md5_value,info,json_task);
    LOG4CXX_INFO(logger_,json_task["MSG"]);
    if (!ret)
        return ERROR;
    json_task["MD5"] = md5_value;
    if (!info.compare(0,7,"dolphin"))
    {
        return MONITOR;
    }
    else
    {
	json json_command;
	try
	{
            json_command = json::parse(info);
	    if (json_command.find("task") == json_command.end())
		throw 1;

            json flag = json_command["task"];
	    if (!flag.is_string())
		throw 1;
            else if (flag == "mongodb")
	    {
                json_task["MSG"] = json_command;
                return MONGODB;
	    }
	    else
		throw 1;
	}
	catch (...)
	{
            return UNKNOWN;
	}
    }
}

void Parse::do_work(std::string& md5_value)
{
    json json_command;
    CommandObject co = parse_command(md5_value,json_command);
    Command* tmp;
    switch (co)
    {
        case ERROR:
            return;
        case UNKNOWN:
        {
            int fd = get_fd(md5_value);
            if (fd != -1)
            {
                std::string er{"unknown format"};
	        write(fd,er.c_str(),er.size());
                remove_task(fd,false);
            }
            return;
        }
        case MONITOR:
            tmp = inst_->create_monitor_inst();
            break;
        case MONGODB:
            tmp = inst_->create_task_inst();
            break;
    }
    int ret = tmp->run_command(json_command);
    if (ret != -1)
    {
        int fd = get_fd(md5_value);
        if (fd != -1)
        {
            tmp->write_result(fd);
            remove_task(fd,false);
        }
    }
    delete tmp;
}

