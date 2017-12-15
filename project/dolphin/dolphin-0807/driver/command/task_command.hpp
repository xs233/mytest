#ifndef __TASK_COMMAND_H__
#define __TASK_COMMAND_H__

#include <bsoncxx/json.hpp>
#include <bsoncxx/builder/basic/array.hpp>
#include <bsoncxx/builder/basic/document.hpp>
#include <bsoncxx/builder/basic/kvp.hpp>
#include <bsoncxx/types.hpp>
#include <mongocxx/stdx.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/instance.hpp>
#include <mongocxx/collection.hpp>
#include "command.hpp"
using bsoncxx::builder::basic::kvp;
using bsoncxx::builder::basic::sub_array;


class TaskCommand : public Command
{
public:
    TaskCommand();
    ~TaskCommand();
    virtual bool run_command(json& json_command) override;
    virtual bool write_result(int fd) override;

private:
    enum {
        TIMESTAMP_ERROR,
        MONGODB_ERROR,
        EMPTY_NUM,
        ASSIGN_ERROR    
    };
    int total_number_;
    std::map<std::string,json> col_info_;
    json json_msg_;
    void set_result(int code,const char* info);
    bool json_find(const char* key);
    bool check_executor();
    bool check_timestamp();
    bool set_col_info(mongocxx::database& db,std::string& col);
    long get_number(mongocxx::collection& coll,long timestamp,bool left_or_right);
};

#endif
