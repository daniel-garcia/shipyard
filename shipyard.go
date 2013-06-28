package shipyard

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// API for Shipyard service
type ShipyardService struct {
	Hosts    map[string]*Host
	Services map[string]*Service
	datafile string
}

// Save the shipyard datafile.
func (server *ShipyardService) save() error {
	data, err := json.Marshal(server)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(server.datafile, data, 0600)
	return err
}

// Create a new Shipyard instance.
func NewShipyardService(datafile string) (server *ShipyardService, err error) {
	server = new(ShipyardService)
	server.Hosts = make(map[string]*Host)
	server.datafile = datafile

	data, err := ioutil.ReadFile(datafile)
	if err == nil {
		err = json.Unmarshal(data, server)
		return server, err
	}
	if os.IsNotExist(err) {
		err = server.save()
	}
	return server, err
}

type ShipyardError struct {
	ErrorMessage string
}

func (e ShipyardError) Error() string {
	return e.ErrorMessage
}

// The Shoreman API client
type ShipyardClient struct {
	connectionString string
	rpcClient        *rpc.Client
	tcpConnection    net.Conn
}

// Create a ShipyardService client.
func NewShipyardClient(connectionString string) (client *ShipyardClient, err error) {
	if len(connectionString) == 0 {
		connectionString = "localhost:4979"
	}
	client = new(ShipyardClient)
	client.connectionString = connectionString
	client.tcpConnection, err = net.Dial("tcp", client.connectionString)
	if err != nil {
		return client, err
	}
	client.rpcClient = jsonrpc.NewClient(client.tcpConnection)
	return client, nil
}

// Call GetHosts() on the ShipyardService
func (client *ShipyardClient) GetHosts() (hosts map[string]*Host, err error) {
	hosts = make(map[string]*Host)
	reply := &GetHostsReply{}
	err = client.rpcClient.Call("ShipyardService.GetHosts", nil, &reply)
	return reply.Hosts, err
}

type GetHostsReply struct {
	Hosts map[string]*Host
}

// Return Host obect to client
func (s *ShipyardService) GetHosts(unused *interface{}, reply *GetHostsReply) error {
	*reply = GetHostsReply{}
	log.Printf("returning %d hosts back to client", len(s.Hosts))
	reply.Hosts = s.Hosts
	return nil
}
