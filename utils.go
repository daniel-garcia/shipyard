package shipyard

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var urandomFilename = "/dev/urandom"

func newUuid() (string, error) {
	f, err := os.Open(urandomFilename)
	if err != nil {
		return "", err
	}
	b := make([]byte, 16)
	defer f.Close()
	f.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, err
}

var hostIdCmd = exec.Command("/usr/bin/hostid")

// hostId retreives the system's unique id, on linux this maps
// to /usr/bin/hostid.
func hostId() (hostid string, err error) {
	cmd := hostIdCmd
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return hostid, err
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	hostid, err = reader.ReadString('\n')
	return strings.TrimSpace(hostid), err
}

// Path to meminfo file. Placed here so getMemorySize() is testable.
var meminfoFile = "/proc/meminfo"

// getMemorySize attempts to get the size of the installed RAM.
func getMemorySize() (size int64, err error) {
	file, err := os.Open(meminfoFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	for err == nil {
		if strings.Contains(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return 0, err
			}
			size, err := strconv.Atoi(parts[1])
			if err != nil {
				return 0, err
			}
			return int64(size * 1024), nil
		}
		line, err = reader.ReadString('\n')
	}
	return 0, err
}
