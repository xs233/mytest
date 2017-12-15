#include <stdlib.h>
#include <iostream>
#include <fstream>
#include "config.hpp"

Config::Config()
{

}

Config::~Config()
{
    if (!ip_slave_.empty())
        ip_slave_.clear();
}

bool Config::start_executor()
{
    for (auto& item : ip_slave_)
    {
        std::string shell_command = "ssh " + item + " \"cd  dolphin/bin; ./executor " + ip_master_ + " " + port_ + "\"";
	std::cout << shell_command;
	int ret = system(shell_command.c_str());
        if (ret == -1)
	{
            std::cerr << "shell command failed!" << std::endl;
            return false;
	}
    }
    return true;
}

bool Config::load_file()
{
    std::ifstream file("node.conf",std::ios::in | std::ios::binary);
    if (!file)
    {
        std::cerr << "open config file failed!" << std::endl;
	return false;
    }

    std::string ip;
    while (file >> ip)
    {
        size_t dot = ip.find(':');
        if (dot != std::string::npos)
	{
	    ip_master_ = ip.substr(0,dot);
	    port_ = ip.substr(dot + 1);
	}
	else
	    ip_slave_.push_back(ip);
    }
    
    for (const auto& it : ip_slave_)
	std::cout << it << std::endl;
    file.close();
    return true;
}

int Config::get_port()
{
    return std::stoi(port_);
}

