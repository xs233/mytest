#include <cstdint>
#include <chrono>
#include <iostream>
#include <fstream>
#include <vector>
#include <bsoncxx/json.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/stdx.hpp>
#include <mongocxx/client.hpp>
#include <mongocxx/instance.hpp>
#include <mongocxx/collection.hpp>

using bsoncxx::builder::stream::close_array;
using bsoncxx::builder::stream::close_document;
using bsoncxx::builder::stream::document;
using bsoncxx::builder::stream::finalize;
using bsoncxx::builder::stream::open_array;
using bsoncxx::builder::stream::open_document;

int main(int,char**)
{
	using namespace std::chrono;
	milliseconds begin_time = duration_cast< milliseconds >(system_clock::now().time_since_epoch());
	mongocxx::instance inst{};
	mongocxx::client client{mongocxx::uri{"mongodb://10.46.215.19:27017"}};
	mongocxx::database db = client["test"];
	mongocxx::collection coll = db["mulidata"];
	mongocxx::cursor cursor = coll.find(document{} << "device" << "66" << finalize);
	std::ofstream file;
	file.open("feature.txt",std::ios::out | std::ios::binary);
	if (!file.is_open())
	{
	    std::cout << "open file failed!" << "\n";
		return -1;
	}
	for(auto doc : cursor) {
	    std::cout << bsoncxx::to_json(doc) << "\n";
	    file << bsoncxx::to_json(doc) << "\n";
	}
	file.close();
	milliseconds end_time = duration_cast< milliseconds >(system_clock::now().time_since_epoch());
	std::cout << "spent time:" << end_time.count() - begin_time.count() << "\n";
    return 0;
}
