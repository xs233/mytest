#include <iostream>
#include "helloworld.pb.h"

google::protobuf::RpcChannel* channel;
google::protobuf::RpcController* controller;
helloworld::Greeter* service;
helloworld::HelloRequest request;
helloworld::HelloReply response;

class MyRpcChannel : public google::protobuf::RpcChannel {
public:
    MyRpcChannel(std::string ip,int port) : ip_(ip), port_(port) {};
    ~MyRpcChannel() {};
    void CallMethod(const google::protobuf::MethodDescriptor* method,
                    google::protobuf::RpcController* controller,
                    const google::protobuf::Message* request,
                    google::protobuf::Message* response,
                    google::protobuf::Closure* done) {
        std::cout << "Hello world!" << '\n';
    };

private:
    std::string ip_;
    int port_;
};

class MyRpcController : public google::protobuf::RpcController {
public:
    MyRpcController() {};
    ~MyRpcController() {};
    void Reset() {};
    bool Failed() const { return true; };
    std::string ErrorText() const { return "failed"; };
    void StartCancel() {};
    void SetFailed(const std::string& reason) {};
    bool IsCanceled() const { return true; };
    void NotifyOnCancel(google::protobuf::Closure* callback) {};
};

void Done() {
    delete service;
    delete channel;
    delete controller;
}

void DoGreeter() {
    channel = new MyRpcChannel("somehost.example.com", 1234);
    controller = new MyRpcController;
    service = new helloworld::Greeter::Stub(channel);

    // Set up the request.
    request.set_name("hello");

    // Execute the RPC.
    service->SayHello(controller, &request, &response, google::protobuf::NewCallback(&Done));
}

int main() {
    DoGreeter();
    return 0;
}