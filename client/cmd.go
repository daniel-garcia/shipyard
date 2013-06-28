
package main

import (
    "github.com/daniel-garcia/shipyard"
    "flag"
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
}
