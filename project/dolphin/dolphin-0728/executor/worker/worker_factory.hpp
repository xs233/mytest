#ifndef __WORKER_FACTORY_H__
#define __WORKER_FACTORY_H__
class MongodbWorker;

class WorkerFactory
{
public:
    WorkerFactory() {};
    virtual ~WorkerFactory() {};
    virtual MongodbWorker* create_mongodb_inst() = 0;

};

#endif
