#include "worker_factory_inst.hpp"
#include "mongodb_worker.hpp"

WorkerFactoryInst::WorkerFactoryInst()
{

}

WorkerFactoryInst::~WorkerFactoryInst()
{

}

MongodbWorker* WorkerFactoryInst::create_mongodb_inst()
{
    return new MongodbWorker("MongodbWorker");
}
