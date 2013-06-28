package main

import (
	"flag"
	"fmt"
	"github.com/daniel-garcia/shipyard"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

var port string
var runShipyardService bool = false
var shipyardDatafile string
var runShoremanService bool = false

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of shoreman: %s [options] run\n", os.Args[0])
	flag.PrintDefaults()
}
func init() {
	flag.StringVar(&port, "port", ":4979", "port to listen on")
	flag.BoolVar(&runShipyardService, "shoreman", true, "run shoreman service")
	flag.BoolVar(&runShoremanService, "shipyard", true, "run shipyard service")
	flag.StringVar(&shipyardDatafile,
		"shipyard-datafile",
		"/tmp/shipyard.json",
		"path to the shipyard data file")
	flag.Usage = usage
}

func startServer() {

	flag.Parse()
	if runShipyardService == false && runShoremanService == false {
		fmt.Fprintf(os.Stderr, "At least one service must be run: -shoreman and/or -shipyard")
	}

	server := rpc.NewServer()
	if runShoremanService {
		log.Print("Starting shoreman service.")
		shoreman := new(shipyard.Shoreman)
		server.Register(shoreman)
	}

	if runShipyardService {
		log.Print("Starting shipyard service.")
		shipyard, err := shipyard.NewShipyardService(shipyardDatafile)
		if err != nil {
			panic(err)
		}
		server.Register(shipyard)
	}

	log.Print("Starting web service.")
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	l, e := net.Listen("tcp", port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	log.Printf("Listening on %s", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	startServer()
}
