#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <iostream>
#include "driver.hpp"
#include "config/config.hpp"

Driver::Driver()
{

}

Driver::~Driver()
{

}

bool Driver::start_driver()
{
    if (!Config::instance().load_file())
    {
        std::cerr << "load config file failed!" << std::endl;
        return false;
    }

    int socket_fd = socket(PF_INET,SOCK_STREAM,0);
    if (socket_fd == -1)
    {
        std::cerr << "socket failed!"  << std::endl;
        return false;
    }

    sockaddr_in addr = {0};
    socklen_t addrlen = sizeof(addr);
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = htonl(INADDR_ANY);
    addr.sin_port = htons(Config::instance().get_port());
    if (bind(socket_fd,(struct sockaddr*)&addr,addrlen))
    {
        std::cerr << "bind failed!" << std::endl;
        close(socket_fd); 
    }

    if (listen(socket_fd,10))
    {
        std::cerr << "listen failed!" << std::endl;
        close(socket_fd); 
    }

    if (!Config::instance().start_executor())
    {
        std::cerr << "start_executor failed!" << std::endl;
        return false;
    }
    close(socket_fd);
}
