#include <iostream>
#include <vector>
#include "../config/config.hpp"
#include "task_command.hpp"

TaskCommand::TaskCommand() : Command("TaskCommand"), total_number_(0)
{

}

TaskCommand::~TaskCommand()
{

}

void TaskCommand::set_result(int code,const char* info)
{
    json error;
    error["task"] = "mongodb";
    error["code"] = code;
    error["error_info"] = info;
    result_ = error.dump();
    LOG4CXX_ERROR(logger_,info);
}

bool TaskCommand::run_command(json& json_command)
{
    std::cout << "json:" << json_command << '\n';
    json_msg_ = json_command["MSG"];
    if (check_executor()) {
        // TODO::
    } else {
        if (!check_timestamp()) {
            set_result(TIMESTAMP_ERROR,"format error!");
            return true;
        }
        mongocxx::instance inst{};
        mongocxx::client client{mongocxx::uri{"mongodb://10.46.215.19:27017"}};
        mongocxx::database db = client["test"];
        std::vector<std::string> device = json_msg_["device"];
        total_number_ = 0;
        for (auto item : device) {
            if (!set_col_info(db,item)) {
                set_result(MONGODB_ERROR,"exception on mongodb!");
                return true;
            }
        }
        if (!total_number_) {
            set_result(EMPTY_NUM,"the total_number is 0!");
            return true;
        }
        std::string md5_value = json_command["MD5"];
        if (!Config::instance().assign_work(md5_value,total_number_,col_info_)) {
            set_result(ASSIGN_ERROR,"assign the work error!");
            return true;
        }
    }
    
    return false;
}

bool TaskCommand::write_result(int fd)
{
    if (write(fd,result_.c_str(),result_.size()) == -1) {
        LOG4CXX_ERROR(logger_,"write monitor fd failed!");
        return false;
    }
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
    if (!(json_find("beginTime") && json_find("endTime") && json_find("device"))) {
        LOG4CXX_ERROR(logger_,"timestamp error!");
        return false;
    }
    return true;
}

bool TaskCommand::set_col_info(mongocxx::database& db,std::string& col)
{
    long left_value,right_value;
    mongocxx::collection coll = db[col];
    long left_time = json_msg_["beginTime"];
    left_value = get_number(coll,left_time,true);
    long right_time = json_msg_["endTime"];
    right_value = get_number(coll,right_time,false);
    left_value = 1498200592015;
    if (left_value && right_value && (left_value <= right_value)) {
        total_number_ += (right_value - left_value + 1);
        json num;
        num["beginTime"] = left_value;
        num["endTime"] = right_value;
        col_info_.insert(std::pair<std::string,json>(col,num));
    }
    return true;
}

long TaskCommand::get_number(mongocxx::collection& coll,long timestamp,bool left_or_right)
{
    bsoncxx::builder::basic::document basic_builder{};
    mongocxx::options::find opts;
    opts.limit(1);
    if (left_or_right) {
        basic_builder.append(kvp("timestamp",
                                        [&timestamp](bsoncxx::builder::basic::sub_document sub_builder)
                                                    {sub_builder.append(kvp("$gte",bsoncxx::types::b_int64{timestamp}));}));
    } else {
        basic_builder.append(kvp("timestamp",
                                        [&timestamp](bsoncxx::builder::basic::sub_document sub_builder)
                                                    {sub_builder.append(kvp("$lte",bsoncxx::types::b_int64{timestamp}));}));
    }
    mongocxx::cursor cursor = coll.find(basic_builder.view(),opts);
    long num = 0;
    for (auto doc : cursor) {
        std::string res = bsoncxx::to_json(doc);
        json json_res = json::parse(res);
        if (json_res.find("timestamp") == json_res.end())
            return 0;
        num = json_res["timestamp"];
    }
    std::cout << num << '\n';
    return num;
}

