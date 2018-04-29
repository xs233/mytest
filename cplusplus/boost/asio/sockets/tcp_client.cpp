#include <iostream>
#include <array>
#include <boost/asio.hpp>

void sync_client() {
    boost::asio::io_service io_service;
    boost::asio::ip::tcp::socket tcp_socket(io_service);
    boost::asio::ip::tcp::resolver resolver(io_service);
    boost::asio::ip::tcp::resolver::results_type endpoints = resolver.resolve("127.0.0.1", "daytime");
    try {
        boost::asio::connect(tcp_socket, endpoints);
        while (true) {
            std::array<char, 1024> buf;
            boost::system::error_code error;
            size_t len = tcp_socket.read_some(boost::asio::buffer(buf), error);
            if (boost::asio::error::eof == error) {
                std::cout << "eof" <<  '\n';
                break;
            } else if (error) {
                throw boost::system::system_error(error);
            }

            std::cout.write(buf.data(), len);
        }
    } catch (std::exception& e) {
        std::cerr << e.what() << '\n';
    }
}

int main() {
    sync_client();
    return 0;
}
