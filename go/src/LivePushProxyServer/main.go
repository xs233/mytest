package main

import (
	"flag"
	"strconv"
	"LivePushProxyServer/server"
)

func main() {
	port := flag.Int("p", 8090, "server listen port")
	flag.Parse()
	server.Server(":" + strconv.Itoa(*port))
}