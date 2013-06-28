
package shipyard

import (
	"os"
	"runtime"
	"time"
	"encoding/json"
	"io/ioutil"
	"fmt"
)




type Network struct {
	Address string
	Netmask string
}

func NewNetwork(address, netmask string) (network *Network, err error) {
	network = new(Network)
	network.Address = address
	network.Netmask = netmask
	return network, err
}
type Host struct {
	Name string
	Id string
	Cores int
	Memory uint64
	PrivateNetwork Network
	LastUpdated time.Time
	controller *ServiceController
}

func (h *Host) String() string {
	return fmt.Sprintf("Host (%s, %s)", h.Name, h.Id)
}

type ResourcePool struct {
	Name string
	Id string
	Cores uint8
	Memory int64
	HostsRefs []string
	controller *ServiceController
}

type Image struct {
	Name string
	Id string
	Tags map[string] string
	controller *ServiceController
}

type TransportType string

const (
	TCP TransportType = "tcp"
	UDP TransportType = "udp"
)

type PortType struct {
	Port uint16
	Transport TransportType
	Application string
	controller *ServiceController
}

type ServicePort struct {
        ServiceId string // unique ID for a Service
        Port PortType   // the Port to map to the Service
	controller *ServiceController
}

type Service struct {
	Name string
	Id string
	ImageRef string
	Description string
	Startup string
	Shutdown string
	Priority int8
	Endpoints map[uint16] *PortType
	ServicePorts []*ServicePort
	controller *ServiceController
}

type ServiceController struct {
	Hosts map[string] *Host
	ResourcePools map[string] *ResourcePool
	SystemId string
	Services map[string] *Service
	filename string
}

// Create  a new ServiceController.
func NewServiceController(filename string) (c *ServiceController, err error) {
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		c = new(ServiceController)
		err = json.Unmarshal(b, &c)
		return c, err
	}
	c = new(ServiceController)
	c.Hosts = make(map[string] *Host)
	c.ResourcePools = make(map[string] *ResourcePool)
	c.Services = make(map[string] *Service)
	c.SystemId = newUuid()
	c.filename = filename
	return c, c.Save()
}


// Create a new Host struct from the running host's values
func CurrentContextAsHost() (host *Host, err error) {
    cpus := runtime.NumCPU()
    memory, err := getMemorySize()
    if err != nil {
        return nil, err
    }
    host = new(Host)
    hostname, err := os.Hostname()
    if err != nil {
        return nil, err
    }
    host.Name = hostname
    hostid_str, err := hostId()
    if err != nil {
        return nil, err
    }
    host.Id = hostid_str
    host.PrivateNetwork = Network{}
    host.Cores = cpus
    host.Memory = memory
    host.LastUpdated = time.Now()
    return host, err
}

// Save the ServiceController to disk.
func (c *ServiceController) Save() error {
	b, err := json.Marshal(c)
	if err != nil { return err }
	return ioutil.WriteFile(c.filename, b, 0600)
}

func (c *ServiceController) AddHost(host *Host) error {
	c.Hosts[host.Id] = host
	return c.Save()
}


