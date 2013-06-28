package main

import (
	"flag"
	"github.com/daniel-garcia/shipyard"
	"log"
)

var port string

func init() {
	flag.StringVar(&port, "port", "localhost:4979", "Port of shoreman service.")
}

func main() {

	client, err := shipyard.NewShoremanClient(port)
	if err != nil {
		panic(err)
	}
	host, err := client.HostInfo()
	log.Print(host.String())

	network, err := client.DockerNetwork()
	log.Print(network)

	err = client.AddRoute("192.168.122.137", shipyard.Network{"172.168.1.0", "255.255.255.0"})
	if err != nil {
		log.Printf(err.Error())
	}

	client.Close()

	shipyardClient, err := shipyard.NewShipyardClient(port)
	if err != nil {
		panic(err)
	}
	hosts, err := shipyardClient.GetHosts()
	if err != nil {
		panic(err)
	}

	log.Printf("Got %d hosts.", len(hosts))

	for _, host = range hosts {
		log.Printf("%s", host)
	}

}
