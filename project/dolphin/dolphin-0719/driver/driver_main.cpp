#include <unistd.h>
#include <syslog.h>
#include <sys/stat.h>
#include <iostream>
#include "driver.hpp"

bool set_daemon()
{
    pid_t pid = fork();
    if (pid == -1)
    {
        std::cerr << "fork failed!" << std::endl;
        exit(EXIT_FAILURE);
    }
    else if (pid == 0)
    {
        if (setsid() < 0)
        {
            std::cerr << "setsid failed!" << std::endl;
            return false;
        }

        umask(0);
        chdir("/");
        for (int x = sysconf(_SC_OPEN_MAX); x>=0; x--)
            close (x);
        openlog ("dolphin_daemon", LOG_PID, LOG_DAEMON);
    }
    else
        exit(EXIT_SUCCESS);
    return true;
}

int main()
{
    if (set_daemon())
    {
        Driver driver;
        syslog (LOG_NOTICE, "driver start");
        driver.start_driver();
    }
    return 0;
}
