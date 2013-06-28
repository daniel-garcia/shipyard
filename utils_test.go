package shipyard

import (
	"os/exec"
	"testing"
)

// Test newUuid()
func TestNewUuid(t *testing.T) {

	urandomFilename = "testfiles/urandom_bytes"

	uuid, err := newUuid()
	if err != nil {
		t.Errorf("Did not expect error: %s", err)
		t.Fail()
	}
	expectedUuid := "1102c395-e94b-0a08-d1e9-307e31a5213e"
	if uuid != expectedUuid {
		t.Errorf("uuid: expected %s, got %s", expectedUuid, uuid)
		t.Fail()
	}
}

// Test getMemorySize()
func TestGetMemorySize(t *testing.T) {

	// alter the file getMemorySize() is looking at
	meminfoFile = "testfiles/meminfo"
	size, err := getMemorySize()
	if err != nil {
		t.Errorf("Failed to parse memory file: %s", err)
		t.Fail()
	}
	expectedSize := int64(33660776448)
	if size != expectedSize {
		t.Errorf("expected %d, received %d ", expectedSize, size)
		t.Fail()
	}
}

// Test hostInfo()
func TestHostInfo(t *testing.T) {

	// alter the command that hostId() is running
	hostIdCmd = exec.Command("cat", "testfiles/hostid")

	hostid, err := hostId()
	if err != nil {
		t.Errorf("Could not retrieve hostid: %s", err)
		t.Fail()
	}
	expectedHostId := "007f0101"
	if hostid != expectedHostId {
		t.Errorf("Hostid, expected: %s, got %s", expectedHostId, hostid)
		t.Fail()
	}
}
