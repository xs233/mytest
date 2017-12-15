#ifndef __TASK_COMMAND_H__
#define __TASK_COMMAND_H__

#include <mongocxx/collection.hpp>
#include "command.hpp"

class TaskCommand : public Command
{
public:
    TaskCommand();
    ~TaskCommand();
    virtual bool run_command(json& json_command) override;
    virtual bool write_result(int fd) override;

private:
    int totoal_number_;
    std::map<std::string,json> col_info_;
    json json_msg_;
    void set_result(int code,std::string info);
    bool json_find(const char* key);
    bool check_executor();
    bool check_timestamp();
    bool set_col_info(mongocxx::database& db,std::string& col);
    int get_number(mongocxx:collection& coll,long timestamp,bool left_or_right);
    bool assign_work(int left_value,int right_value);
};

#endif
