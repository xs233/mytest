.PHONY:clean all
CPPFLAGS=-g -Wall
SRC = $(wildcard driver/*.cpp)
#BIN = $(SRC:%.cpp=%)
BIN = bin/driver
SRC1 = $(wildcard driver/config/*.cpp)
SRC2 = $(wildcard driver/command/*.cpp)
SRC3 = $(wildcard executor/worker/*.cpp)

all:$(BIN)
$(BIN): $(SRC) $(SRC1) $(SRC2)
	g++ $(CPPFLAGS) -std=c++11 $^ -o $@ -pthread -llog4cxx -lcryptopp
	g++ $(CPPFLAGS) -std=c++11 $(SRC3) executor/executor.cpp executor/executor_main.cpp -o bin/executor -pthread -llog4cxx

clean:
	rm  bin/driver bin/executor
