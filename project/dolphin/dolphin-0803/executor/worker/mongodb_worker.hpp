#ifndef __MONGODB_WORKER_H__
#define __MONGODB_WORKER_H__
#include "worker.hpp"

class MongodbWorker : public Worker
{
public:
    MongodbWorker(const char* log_name);
    virtual  ~MongodbWorker() override;
    virtual void run(int fd,json& task_json) override;

};

#endif
