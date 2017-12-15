.PHONY:clean all
CPPFLAGS=-g -Wall
SRC = $(wildcard driver/*.cpp)
#BIN = $(SRC:%.cpp=%)
BIN = bin/driver
SRC1 = $(wildcard driver/config/*.cpp)
SRC2 = $(wildcard driver/command/*.cpp)

all:$(BIN)
$(BIN): $(SRC) $(SRC1) $(SRC2)
	g++ $(CPPFLAGS) -std=c++11 $^ -o $@ -pthread -llog4cxx
	g++ $(CPPFLAGS) -std=c++11 executor/executor.cpp -o bin/executor

clean:
	rm  bin/driver bin/executor
