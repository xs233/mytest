package main

import (
	"flag"
	"log"
	"net/http"
	"tinyecho/core"

	"golang.org/x/net/websocket"
)

var (
	pPtr    = flag.String("p", "", "listen port")
	portPtr = flag.String("port", "", "listen port")
)

func main() {
	flag.Parse()

	port := ""

	switch true {
	case *pPtr != "":
		port = *pPtr
		port = *pPtr
	case *portPtr != "":
		port = *portPtr
	default:
		port = "9000"
	}
    
    http.HandleFunc("/echo",
        func (w http.ResponseWriter, req *http.Request) {
            s := websocket.Server{Handler: websocket.Handler(core.ServeWS)}
            s.ServeHTTP(w, req)
        })
        
	// http.Handle("/echo", websocket.Server{Handler : core.ServeWS})
	http.Handle("/", websocket.Handler(core.ServeWS))
	http.HandleFunc("/online", core.QueryOnlineClients)
	http.HandleFunc("/publish", core.PublishMessage)
	log.Printf("Listening and Serving on :%v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
