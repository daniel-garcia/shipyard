
package main

import (
    "flag"
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "github.com/daniel-garcia/shipyard"
)

var port string

func init() {
	flag.StringVar(&port, "port", ":4979", "port to listen on")
}

func startServer() {
	flag.Parse()

	shoreman := new(shipyard.Shoreman)

	server := rpc.NewServer()
	server.Register(shoreman)

	log.Print("Starting web service.")
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	l, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	log.Printf("Listening on %s", port)
	for {
		conn, err := l.Accept()
		if err != nil { log.Fatal(err) }
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}


func main() {
	startServer()
}

