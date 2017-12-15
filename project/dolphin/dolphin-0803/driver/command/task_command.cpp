#include <iostream>
#include <vector>
#include <bsoncxx/json.hpp>
#include <bsoncxx/builder/basic/array.hpp>
#include <bsoncxx/builder/basic/document.hpp>
#include <bsoncxx/builder/basic/kvp.hpp>
#include <bsoncxx/types.hpp>
#include <mongocxx/stdx.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/instance.hpp>
#include <mongocxx/collection.hpp>
#include "task_command.hpp"
using bsoncxx::builder::basic::kvp;
using bsoncxx::builder::basic::sub_array;

TaskCommand::TaskCommand() : Command("TaskCommand"), total_number_(0)
{

}

TaskCommand::~TaskCommand()
{

}

void set_result(int code,std::string info)
{
    json error["task"] = "mongodb";
    error["code"] = code;
    error["error_info"] = info;
    result_ = error;
    LOG4CXX_ERROR(logger_,info);
}

bool TaskCommand::run_command(json& json_command)
{
    std::cout << "json:" << json_command << '\n';
    json_msg_ = json_command;
    if (check_executor()) {
        // TODO::
    } else {
        if (!check_timestamp()) {
            set_result(1,"format error!");
            return true;
        }
        mongocxx::instance inst{};
        mongocxx::client client{mongocxx::uri{"mongodb://10.46.215.19:27017"}};
        mongocxx::database db = client["test"];
        std::vector<std::string> device = json_msg_["device"];
        for (auto item : device) {
            if (!set_col_info(db,item)) {
                set_result(2,"exception on mongodb!");
                return true;
            }
        }
    }
    
    return false;
}

bool TaskCommand::write_result(int fd)
{
    return true;
}

bool TaskCommand::json_find(const char* key)
{
    auto ret = json_msg_.find(key);
    if (ret == json_msg_.end())
        return false;
    return true;
}

bool TaskCommand::check_executor()
{
    if (json_find("executor"))
        return true;
    return false;
}

bool TaskCommand::check_timestamp()
{
    if (json_find("beginTime") && json_find("endTime") && json_find("device")) {
        LOG4CXX_ERROR(logger_,"timesatmp error!");
        return false;
    }
    return true;
}

bool TaskCommand::set_col_info(mongocxx::database& db,std::string& col)
{
    int left_value,right_value;
    mongocxx::collection coll = db[col];
    try {
        long left_time = json_msg_["beginTime"];
        left_value = get_number(coll,left_time,true);
        long right_time = json_msg_["endTime"];
        right_value = get_number(coll,right_time,false);
    } catch (mongocxx::logic_error ) {
        LOG4CXX_ERROR(logger_,"coll.find error!");
        return false;
    }
    if (left_value && right_value && (left_value <= right_value)) {
        total_number_ += (right_value - left_value + 1);
        json num["beginNum"] = left_value;
        num["endNum"] = right_value;
        col_info_.insert(std::pair<std::string,json>(col,num));
    }
    return true;
}

int TaskCommand::get_number(mongocxx:collection& coll,long timestamp,bool left_or_right)
{
    bsoncxx::builder::basic::document basic_builder{};
    mongocxx::options::find opts;
    opts.limit(1);
    if (left_or_right) {
        basic_builder.append(kvp("timestamp",
                                        [](bsoncxx::builder::basic::sub_document sub_builder)
                                                    {sub_builder.append(kvp("$gte",bsoncxx::types::b_int64{timestamp}));}));
    } else {
        basic_builder.append(kvp("timestamp",
                                        [](bsoncxx::builder::basic::sub_document sub_builder)
                                                    {sub_builder.append(kvp("$lte",bsoncxx::types::b_int64{timestamp}));}));
    }
    mongocxx::cursor cursor = coll.find(basic_builder.view(),opts);
    int num = 0;
    for (auto doc : cursor) {
        std::string res = bsoncxx::to_json(doc);
        json json_res = json::parse(res);
        if (json_res.find("number") == json_res.end())
            return 0;
        num = json_res["number"];
    }
    return num;
}

bool TaskCommand::assign_work(int left_value,int right_value)
{
    int remaining_cores = Config::instance().get_total_cores();
    if (!remaining_cores) {
        set_result(3,"there is no vaild cores!");
        return false;
    }
    int per_core_task = (total_number_ - total_number_ % remaining_cores) / remaining_cores;

}

