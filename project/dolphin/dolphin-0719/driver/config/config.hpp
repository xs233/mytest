#ifndef __CONFIG_HPP__
#define __CONFIG_HPP__

#include <string>
#include <vector>

class Config
{
public:
    static Config& instance()
    {
        static Config inst;
        return inst;
    };
    bool load_file();
    bool start_executor();
    int get_port();
    Config(const Config& ) = delete;
    Config& operator=(const Config& ) = delete;

private:
    Config();
    ~Config();
    std::vector<std::string> ip_slave_;
    std::string ip_master_;
    std::string port_;

};

#endif
