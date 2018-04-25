#include <ctime>
#include <string>
#include <iostream>
#include <boost/asio.hpp>

using boost::asio::ip::tcp;

std::string make_daytime_string() {
    time_t now = time(0);
    return ctime(&now);
}

void sync_server() {
    boost::asio::io_service io_service;
    tcp::acceptor acceptor(io_service, tcp::endpoint(tcp::v4(), 13));
    try {
        while (true) {
            tcp::socket socket(io_service);
            acceptor.accept(socket);
            std::string msg = make_daytime_string();
            boost::asio::write(socket, boost::asio::buffer(msg));
        }
    } catch (std::exception& e) {
        std::cout << e.what() << '\n';
    }
}

int main() {
    sync_server();
    return 0;
}