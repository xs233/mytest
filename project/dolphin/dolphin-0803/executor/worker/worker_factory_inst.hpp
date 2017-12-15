#ifndef __WORKER_FACTORY_INST_H__
#define __WORKER_FACTORY_INST_H__
#include "worker_factory.hpp"

class WorkerFactoryInst : public WorkerFactory
{
public:
    WorkerFactoryInst();
    virtual ~WorkerFactoryInst() override;
    virtual MongodbWorker* create_mongodb_inst() override;

};

#endif
