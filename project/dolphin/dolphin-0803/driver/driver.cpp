#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <cstring>
#include "driver.hpp"
#include "config/config.hpp"
#include "command/parse.hpp"

Driver::Driver() : logger_(LoggerPtr(Logger::getLogger("driver"))),recv_buf_(new char[BUF_SIZE])
{
    memset(evs_,'\0',CONN_MAX * sizeof(epoll_event));
}

Driver::~Driver()
{
    if (recv_buf_)
    {
        delete[] recv_buf_;
        recv_buf_ = nullptr;
    }
}

bool Driver::start_driver()
{
    PropertyConfigurator::configure("driver.properties");
    Parse::instance().init();
    
    if (!init_tcp_env())
        return false;
    
    if (!Config::instance().start_executor())
    {
        LOG4CXX_ERROR(logger_, "start_executor failed!");
        close(socket_fd_);
        return false;
    }
    epoll();
    return true;
}

bool Driver::init_tcp_env()
{
    if (!Config::instance().load_file())
    {
        LOG4CXX_ERROR(logger_, "load config file failed!");
        return false;
    }

    socket_fd_ = socket(PF_INET,SOCK_STREAM,0);
    if (socket_fd_ == -1)
    {
        LOG4CXX_ERROR(logger_, "socket failed!");
        return false;
    }

    sockaddr_in addr = {0};
    socklen_t addrlen = sizeof(addr);
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = htonl(INADDR_ANY);
    addr.sin_port = htons(Config::instance().get_port());
    if (bind(socket_fd_,(struct sockaddr*)&addr,addrlen))
    {
        LOG4CXX_ERROR(logger_, "bind failed!");      
        close(socket_fd_);
        return false;
    }

    if (listen(socket_fd_,10))
    {
        LOG4CXX_ERROR(logger_, "listen failed!");
        close(socket_fd_); 
        return false;
    }
    return true;
}

bool Driver::epoll()
{
	epfd_ = epoll_create(1);
	if (epfd_ < 0)
	{
                LOG4CXX_ERROR(logger_, "epoll_create failed!");
		close(socket_fd_);
		return false;
	}

	epoll_event ev;
	ev.events = EPOLLIN;
	ev.data.fd = socket_fd_;
	epoll_ctl(epfd_,EPOLL_CTL_ADD,socket_fd_,&ev);
        LOG4CXX_INFO(logger_, "epoll server start OK...");
	
	int ret = 0;
	while(1)
	{
		ret = epoll_wait(epfd_,evs_,CONN_MAX,-1);
                LOG4CXX_INFO(logger_,"epoll_wait return...");
		if (ret < 0)
		{
			if (EINTR == errno)
			{
                                LOG4CXX_INFO(logger_,"catch EINTR");
				return true;
			}
                        LOG4CXX_ERROR(logger_,"epoll_wait failed!");
                        close(epfd_);
			close(socket_fd_);
			return false;
		} 
		else
		{
			if (run_command(ret) != 0)
                        {
                            close(epfd_);
                            close(socket_fd_);
                            return true;    
                        }
		}
	}
	return true;
}

int Driver::run_command(int cnt)
{
	for (int i=0; i<cnt; ++i)
	{
		if (socket_fd_ == evs_[i].data.fd)
		{
                        LOG4CXX_INFO(logger_,"new client connect");
			sockaddr_in useraddr = {0};
			socklen_t nLen = sizeof(sockaddr_in);
			int nUserFd = accept(socket_fd_,(struct sockaddr*)&useraddr,&nLen);
			if (nUserFd < 0)
			{
                                LOG4CXX_ERROR(logger_,"accept new connect failed,why it happened?");
				continue;
			}
			else
			{
				epoll_event ev = {0};
				ev.events = EPOLLIN;
				ev.data.fd = nUserFd;
				epoll_ctl(epfd_,EPOLL_CTL_ADD,nUserFd,&ev);
				LOG4CXX_INFO(logger_,"new user from " << inet_ntoa(useraddr.sin_addr) << ':' << ntohs(useraddr.sin_port));
				write(nUserFd,"welcome",7);
                                std::string ip{inet_ntoa(useraddr.sin_addr)};
                                Config::instance().insert_client(nUserFd,ip);
			}
		}
		else
		{
                        LOG4CXX_INFO(logger_,"client message process");
                        memset(recv_buf_,'\0',BUF_SIZE);
			int ret = read(evs_[i].data.fd,recv_buf_,BUF_SIZE);
			if (ret < 0)
			{
                                LOG4CXX_ERROR(logger_,"recv msg failed,why it happened?");
				continue;
			}
			else if (0 == ret)
			{
				epoll_ctl(epfd_,EPOLL_CTL_DEL,evs_[i].data.fd,evs_ + i);
				close(evs_[i].data.fd);
                                LOG4CXX_INFO(logger_,"disconnect: " << evs_[i].data.fd);
                                Config::instance().erase_client(evs_[i].data.fd);
                                Parse::instance().remove_task(evs_[i].data.fd,true);
			}
			else
			{
                                //std::string msg{recv_buf_};
                                if (!std::strcmp(recv_buf_,"dolphin exit"))
                                {
                                    LOG4CXX_INFO(logger_,"recv exit msg,epoll over!");
                                    return 1;
                                }
                                //msg = std::to_string(evs_[i].data.fd) + "||" + msg;
                                Parse::instance().push(evs_[i].data.fd,recv_buf_,std::strlen(recv_buf_));
                                LOG4CXX_INFO(logger_,"push msg");
			}
		}
	}
	return 0;
}
