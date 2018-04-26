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

class tcp_connection : public std::enable_shared_from_this<tcp_connection> {
public:
    ~tcp_connection() {
        std::cout << "~tcp_connection()" << '\n';
    }

    static std::shared_ptr<tcp_connection> create(boost::asio::io_service& io_service) {
        return std::shared_ptr<tcp_connection>(new tcp_connection(io_service));
    }

    void start() {
        boost::asio::async_write(socket_, boost::asio::buffer(make_daytime_string()), 
            std::bind(&tcp_connection::write_handler, shared_from_this(), std::placeholders::_1, std::placeholders::_2));
    }

    tcp::socket& get_socket() {
        return socket_;
    }

private:
    tcp_connection(boost::asio::io_service& io_service) : socket_(io_service) {}
    
    void write_handler(const boost::system::error_code& ec, std::size_t bytes_transferred) {
        std::cout << ec << " " << bytes_transferred << '\n';
    }

    tcp::socket socket_;
};

class tcp_server {
public:
    tcp_server(boost::asio::io_service& io_service) : acceptor_(io_service, tcp::endpoint(tcp::v4(), 13)) {
        start();
    }
    ~tcp_server() {}

private:
    void start() {
        auto new_connection = tcp_connection::create(acceptor_.get_executor().context());
        acceptor_.async_accept(new_connection->get_socket(), std::bind(&tcp_server::accept_handler, this, new_connection, std::placeholders::_1));
    }

    void accept_handler(std::shared_ptr<tcp_connection> new_connection, const boost::system::error_code& error) {
        if (error) {
            std::cout << error << '\n';
        } else {
            new_connection->start();
        }
        start();
    }

    tcp::acceptor acceptor_;

};

void async_server() {
    boost::asio::io_service io_service;
    tcp_server server(io_service);
    io_service.run();
}

int main() {
    // sync_server();
    async_server();
    return 0;
}