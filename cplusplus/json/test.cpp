#include <iostream>
#include <sstream>
#include <json.hpp>
#include <vector>
#include <map>

using json = nlohmann::json;


int main()
{    
    json j2 = {
        {"pi", 3141},
        {"happy", true},
        {"name", "Niels"},
        {"nothing", nullptr},
        {"answer", {
            {"everything", 42}
        }},
        {"list", {1, 0, 2}},
            {"object", {
            {"currency", "USD"},
            {"value", 42.99}
        }}
       }; 
    std::cout << j2.dump() << std::endl << j2["object"]["value"] << std::endl << j2["list"][2] << std::endl;
    if (j2.find("aa") == j2.end())
        std::cout << "not find" << '\n';
    if (j2["name"] == "Niels")
        std::cout << j2["aa"] << '\n';
    std::vector<int> ivec = j2["list"];    
    for (auto n : ivec)
        std::cout << n << '\n';
    std::string str;
    int a = j2["pi"];
    std::cout << a << '\n';
    
    json j3;
    j3.push_back("aa");
    j3.push_back("bb");
    std::cout << j3 << '\n';

    std::map<int,int> imap;
    imap.insert(std::pair<int,int>(1,1));
    imap.insert(std::pair<int,int>(2,2));
    imap.insert(std::pair<int,int>(3,3));
    imap.insert(std::pair<int,int>(4,4));
    imap.insert(std::pair<int,int>(5,5));
    for (auto it = imap.begin(); it != imap.end();) {
        if (it->first == 2) {
            // auto tmp = it;
            // ++it;
            imap.erase(it++);
        } else {
            std::cout << "first:" << it->first << '\n';
            int& n = it->second;
            n = 11;
            ++it;
        }
    }
    for (auto it : imap) {
        std::cout << "second:" << it.second << '\n';
    }
    return 0;
}
