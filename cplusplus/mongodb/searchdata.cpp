#include <cstdint>
#include <chrono>
#include <iostream>
#include <fstream>
#include <vector>
#include <bsoncxx/json.hpp>
#include <bsoncxx/builder/basic/array.hpp>
#include <bsoncxx/builder/basic/document.hpp>
#include <bsoncxx/builder/basic/kvp.hpp>
#include <bsoncxx/types.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/stdx.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/instance.hpp>
#include <mongocxx/collection.hpp>

using bsoncxx::builder::basic::kvp;
using bsoncxx::builder::basic::sub_array;

int main(int,char**)
{
	using namespace std::chrono;
	milliseconds begin_time = duration_cast< milliseconds >(system_clock::now().time_since_epoch());
	mongocxx::instance inst{};
	mongocxx::client client{mongocxx::uri{"mongodb://10.46.215.19:27017"}};
	mongocxx::database db = client["test"];
	mongocxx::collection coll = db["litdata"];
	bsoncxx::builder::basic::document basic_builder{};
	mongocxx::options::find opts;
	opts.limit(15);
	basic_builder.append(kvp("timestamp",
				             [](bsoncxx::builder::basic::sub_document sub_builder)
							 {sub_builder.append(kvp("$gt",bsoncxx::types::b_int64{1998201111121}));}));
	basic_builder.append(kvp("device", "66"));
	mongocxx::cursor cursor = coll.find(basic_builder.view(),opts);
	std::cout << bsoncxx::to_json(basic_builder.view()) << "\n";
	// std::cout << coll.count(basic_builder.view()) << "\n";
	
	for(auto doc : cursor) {
	    std::cout << bsoncxx::to_json(doc) << "\n";
	}
    
	/*
	std::ofstream file;
	file.open("feature.txt",std::ios::out | std::ios::binary);
	if (!file.is_open())
	{
	    std::cout << "open file failed!" << "\n";
		return -1;
	}
	for(auto doc : cursor) {
	    // std::cout << bsoncxx::to_json(doc) << "\n";
		file << bsoncxx::to_json(doc) << "\n";
	}
	file.close();
	*/
	milliseconds end_time = duration_cast< milliseconds >(system_clock::now().time_since_epoch());
	std::cout << "spent time:" << end_time.count() - begin_time.count() << "\n";
    return 0;
}
