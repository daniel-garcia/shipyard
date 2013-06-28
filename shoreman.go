
package shipyard

import (
    "net"
    "net/rpc/jsonrpc"
    "os/exec"
	"strings"
)


// API for Shoreman daemon.
type Shoreman struct {}

type ShoremanError struct {
	ErrorMessage string
}

func (e ShoremanError) Error() string {
	return e.ErrorMessage
}
// The Shoreman API client
type ShoremanClient struct {
	connectionString string
}

func NewShoremanClient(connectionString string) (client *ShoremanClient, err error) {
	if len(connectionString) == 0 {
		connectionString = "localhost:4979"
	}
	client = &ShoremanClient{connectionString}

	//conn, err := net.Dial("tcp", connectionString)
	//if err != nil {
//		return client, err
//	}
	//defer conn.Close()
	return client, nil
}

func (client *ShoremanClient) HostInfo() (host *Host, err error) {
	conn, err := net.Dial("tcp", client.connectionString)
	if err != nil {
		return host, err
	}
	//defer conn.Close()

	host = &Host{}
	c := jsonrpc.NewClient(conn)
	err = c.Call("Shoreman.HostInfo", nil, &host)
	return host, err
}

// Return Host obect to client
func (s *Shoreman) HostInfo(unused *interface{}, host *Host) error {
	hostResponse, err := CurrentContextAsHost()
	*host = *hostResponse
	return err
}

func (s *Shoreman) DockerNetwork(unused *interface{}, network *Network) error {
	out, err := exec.Command("route").Output()
	if err != nil { return err }

	for _, line := range(strings.Split(string(out), "\n")) {
		if strings.Contains(line, "docker0") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				return ShoremanError{"Could not parse route output"}
			}
			networkIp := fields[0]
			netmask := fields[2]
			*network = Network{networkIp,netmask}
			return nil
		}
	}
	return ShoremanError{"Did not find docker network"}
}

func (client *ShoremanClient) DockerNetwork() (network *Network, err error) {
	conn, err := net.Dial("tcp", client.connectionString)
	if err != nil {
		return network, err
	}
	//defer conn.Close()

	network = &Network{}
	c := jsonrpc.NewClient(conn)
	err = c.Call("Shoreman.DockerNetwork", nil, &network)
	return network, err
}


