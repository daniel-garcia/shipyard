
package shipyard

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"os/exec"
)

func newUuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	defer f.Close()
	f.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// hostId retreives the system's unique id, on linux this maps
// to /usr/bin/hostid.
func hostId() (hostid string, err error) {
    cmd := exec.Command("/usr/bin/hostid")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return hostid, err
    }
    cmd.Start()
    reader := bufio.NewReader(stdout)
    hostid, err = reader.ReadString('\n')
    return strings.TrimSpace(hostid), err
}

// getMemorySize attempts to get the size of the installed RAM.
func getMemorySize() (size uint64, err error) {
    file, err := os.Open("/proc/meminfo")
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
            return uint64(size * 1024), nil
        }
        line, err = reader.ReadString('\n')
    }
    return 0, err
}

