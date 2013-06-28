package shipyard

import (
	"fmt"
)

// Represents a set of available resources and their contraints.
type ResourcePool struct {
	Name      string
	Id        string
	Cores     int
	Memory    int64
	HostsRefs []string
	service   *ShipyardService
}

// Prettyprint ResourcePool
func (pool *ResourcePool) Save() string {
	return fmt.Sprintf("ResourcePool[%s, %s]", pool.Name, pool.Id)
}

// Get the total amount of cores in the resource pool
func (pool *ResourcePool) GetTotalCores() (cores int, err error) {
	if len(pool.HostsRefs) == 0 {
		return -1, ShipyardError{"Resource pool has no hosts"}
	}
	totalCores := 0
	for _, hostRef := range pool.HostsRefs {
		host, ok := pool.service.Hosts[hostRef]
		if !ok {
			msg := fmt.Sprintf("Data integrity error, could not find hostref "+
				"%s in resource pool %s", hostRef, pool.Id)
			panic(ShipyardError{msg})
		}
		totalCores += host.Cores
	}
	if pool.Cores != 0 || pool.Cores < totalCores {
		return pool.Cores, nil
	}
	return totalCores, nil
}

// Get the total amount of memory in this resource pool.
func (pool *ResourcePool) GetTotalMemory() (totalMemory int64, err error) {

	if len(pool.HostsRefs) == 0 {
		return -1, ShipyardError{"Resource pool has no hosts"}
	}
	for _, hostRef := range pool.HostsRefs {
		host, ok := pool.service.Hosts[hostRef]
		if !ok {
			panic(ShipyardError{"Data integrity error, could not find host"})
		}
		totalMemory += host.Memory
	}
	if pool.Memory != 0 || pool.Memory < totalMemory {
		return pool.Memory, nil
	}
	return totalMemory, nil
}
