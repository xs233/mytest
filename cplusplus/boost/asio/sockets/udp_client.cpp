#include <iostream>
#include <array>
#include <boost/asio.hpp>

void sync_client() {
    boost::asio::io_service io_service;
    boost::asio::ip::udp::socket udp_socket(io_service);
    boost::asio::ip::udp::resolver resolver(io_service);
    boost::asio::ip::udp::endpoint receiver_endpoint = *resolver.resolve(boost::asio::ip::udp::v4(), 
        "127.0.0.1", "daytime").begin();
    try {
        udp_socket.open(boost::asio::ip::udp::v4());
        std::array<char, 1> send_buf  = {{ 0 }};
        udp_socket.send_to(boost::asio::buffer(send_buf), receiver_endpoint);
        std::array<char, 128> recv_buf;
        boost::asio::ip::udp::endpoint sender_endpoint;
        size_t len = udp_socket.receive_from(boost::asio::buffer(recv_buf), sender_endpoint);
        std::cout.write(recv_buf.data(), len);
    } catch (std::exception& e) {
        std::cerr << e.what() << '\n';
    }
}

int main() {
    sync_client();
    return 0;
}