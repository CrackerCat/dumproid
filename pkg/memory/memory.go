package memory

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type memoryMap struct {
	beginAddr int64
	endAddr   int64
	path      string
}

var splitSize = 0x50000000
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, splitSize)
	},
}

func DisplayMemoryMap(pid string, filter string) error {
	mapsPath := fmt.Sprintf("/proc/%s/maps", pid)
	file, err := os.Open(mapsPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		meminfo := strings.Fields(line)
		permission := meminfo[1]
		if permission == filter {
			fmt.Println(scanner.Text())
		}
	}
	return nil
}

func filterMemoryMap(mapsPath string, filter string) ([]memoryMap, error) {
	file, err := os.OpenFile(mapsPath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	memoryMaps := []memoryMap{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		meminfo := strings.Fields(line)
		addrRange := meminfo[0]
		permission := meminfo[1]
		if filter != "" && permission == filter {
			addrs := strings.Split(addrRange, "-")
			beginAddr, _ := strconv.ParseInt(addrs[0], 16, 64)
			endAddr, _ := strconv.ParseInt(addrs[1], 16, 64)
			path := ""
			if len(meminfo) > 5 {
				path = strings.Join(meminfo[5:], " ")
			}
			memoryMaps = append(memoryMaps, memoryMap{
				beginAddr: beginAddr,
				endAddr:   endAddr,
				path:      path,
			})
		}
	}
	return memoryMaps, nil
}

func DisplayMemoryBytes(pid string, beginAddress int64, size int64) error {
	memPath := fmt.Sprintf("/proc/%s/mem", pid)
	memFile, _ := os.Open(memPath)
	defer memFile.Close()

	saved := make([]byte, size)
	memory := readMemory(memFile, saved, beginAddress, size)
	fmt.Printf("%s", hex.Dump(memory))
	return nil
}

func DumpToFile(pid string, permission string, outputPath string) error {
	t := time.Now()
	outputDir := filepath.Join(outputPath, t.Format("20060102150405"))
	if err := os.Mkdir(outputDir, 0777); err != nil {
		return err
	}
	fmt.Printf("Output Dir: %s\n", outputDir)

	mapsPath := fmt.Sprintf("/proc/%s/maps", pid)
	memoryMaps, err := filterMemoryMap(mapsPath, permission)
	if err != nil {
		return err
	}

	memPath := fmt.Sprintf("/proc/%s/mem", pid)
	memFile, _ := os.Open(memPath)
	defer memFile.Close()

	for _, v := range memoryMaps {
		memSize := v.endAddr - v.beginAddr
		dumpFileName := fmt.Sprintf("%x-%x_%s", v.beginAddr, v.endAddr, v.path)
		dumpFileName = strings.Replace(dumpFileName, "/", "_", -1)
		dumpFileName = strings.Replace(dumpFileName, " ", "_", -1)
		dumpPath := filepath.Join(outputDir, dumpFileName)
		fmt.Printf("  Dump File: %s\n", dumpFileName)

		for i := 0; i < (int(memSize)/splitSize)+1; i++ {
			splitIndex := int64((i + 1) * splitSize)
			splittedBeginAddr := v.beginAddr + int64(i*splitSize)
			splittedEndAddr := v.endAddr
			if splitIndex < memSize {
				splittedEndAddr = v.beginAddr + splitIndex
			}
			splittedMemSize := (splittedEndAddr - splittedBeginAddr)
			b := bufferPool.Get().([]byte)[:splittedMemSize]
			memory := readMemory(memFile, b, splittedBeginAddr, splittedMemSize)

			dumpFile, err := os.OpenFile(dumpPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println(err)
				return err
			}
			defer dumpFile.Close()
			dumpFile.Write(memory)
			bufferPool.Put(b)
		}
	}
	return nil
}

func readMemory(memFile *os.File, buffer []byte, beginAddress int64, size int64) []byte {
	r := io.NewSectionReader(memFile, beginAddress, size)
	r.Read(buffer)
	return buffer
}
