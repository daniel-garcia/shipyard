package shipyard

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
	"strings"
)

// API for Shoreman daemon.
type Shoreman struct{}

// An error type for Shoreman operations
type ShoremanError struct {
	ErrorMessage string
}

func (e ShoremanError) Error() string {
	return e.ErrorMessage
}

// The Shoreman API client
type ShoremanClient struct {
	connectionString string
	rpcClient        *rpc.Client
	tcpConnection    net.Conn
}

func (client *ShoremanClient) Close() error {
	err := client.rpcClient.Close()
	if err != nil {
		return err
	}
	return client.tcpConnection.Close()
}

// Create a new Shoreman RPC client
func NewShoremanClient(connectionString string) (client *ShoremanClient, err error) {
	if len(connectionString) == 0 {
		connectionString = "localhost:4979"
	}
	client = new(ShoremanClient)
	client.connectionString = connectionString
	client.tcpConnection, err = net.Dial("tcp", client.connectionString)
	if err != nil {
		return client, err
	}
	client.rpcClient = jsonrpc.NewClient(client.tcpConnection)
	return client, nil
}

// call Shoreman.HostInfo() RPC method
func (client *ShoremanClient) HostInfo() (host *Host, err error) {
	host = &Host{}
	err = client.rpcClient.Call("Shoreman.HostInfo", nil, &host)
	return host, err
}

// Return Host obect to client
func (s *Shoreman) HostInfo(unused *interface{}, host *Host) error {
	hostResponse, err := CurrentContextAsHost()
	*host = *hostResponse
	return err
}

// Return the private network used by docker running on this host.
func (s *Shoreman) DockerNetwork(unused *interface{}, network *Network) error {
	out, err := exec.Command("route").Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "docker0") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				return ShoremanError{"Could not parse route output"}
			}
			networkIp := fields[0]
			netmask := fields[2]
			*network = Network{networkIp, netmask}
			return nil
		}
	}
	return ShoremanError{"Did not find docker network"}
}

// call the Shoreman.DockerNetwork() RPC
func (client *ShoremanClient) DockerNetwork() (network *Network, err error) {
	err = client.rpcClient.Call("Shoreman.DockerNetwork", nil, &network)
	return network, err
}

// Argument to Shoreman.AddRoute()
type AddRouteRequest struct {
	Gateway        string
	NetworkAddress Network
}

// Attempt to add a static route to the Shoreman host
func (s *Shoreman) AddRoute(request AddRouteRequest, unused *interface{}) error {

	// is the network reachable
	err := exec.Command("ping", request.Gateway, "-c", "1").Run()
	if err != nil {
		return ShoremanError{fmt.Sprintf("Could not ping %s", request.Gateway)}
	}

	err = exec.Command("sudo", "route", "add", "-net", request.NetworkAddress.Address, "netmask", request.NetworkAddress.Netmask, "gw", request.Gateway).Run()
	if err != nil {
		return ShoremanError{"Could not add route"}
	}
	return nil
}

// Attempt to add a static route to the Shoreman host
func (client *ShoremanClient) AddRoute(gateway string, network Network) error {
	request := AddRouteRequest{gateway, network}
	return client.rpcClient.Call("Shoreman.AddRoute", request, &struct{}{})
}
