package memory

import (
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestReadMemory(t *testing.T) {
	memFile, _ := os.Open("testdata/proc_test_mem")
	defer memFile.Close()
	size := 5
	saved := make([]byte, size)
	actual := readMemory(memFile, saved, 0x3, int64(size)) // Is it really zero origin?
	expected := []byte{0x3, 0x4, 0x5, 0x6, 0x7}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got memory bytes: %v\nexpected memory bytes: %v", actual, expected)
	}
}

func TestFilterMemoryMap(t *testing.T) {
	memoryMaps, _ := filterMemoryMap("testdata/proc_test_maps", "rw-p")
	expectedAddrs := []string{"12c00000", "6f181000", "6f3bc000", "6f4c5000", "7e0c89c000"}
	expectedPaths := []string{"/dev/ashmem/dalvik-main space (region space) (deleted)",
		"/data/dalvik-cache/arm/system@framework@boot.art",
		"/data/dalvik-cache/arm/system@framework@boot-core-libart.art",
		"/data/dalvik-cache/arm/system@framework@boot-conscrypt.art",
		"",
	}
	for i, v := range memoryMaps {
		expectedAddr, _ := strconv.ParseInt(expectedAddrs[i], 16, 64)
		if v.beginAddr != expectedAddr {
			t.Errorf("got begin addr: %v\nexpected begin addr: %v", v.beginAddr, expectedAddr)
		}
		if v.path != expectedPaths[i] {
			t.Errorf("got begin addr: %v\nexpected begin addr: %v", v.path, expectedPaths[i])
		}
	}
}
