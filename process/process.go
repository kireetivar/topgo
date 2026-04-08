package process

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	PID  int64
	Name string
	CPU  float64
	Mem  float64
}

func GetProcessList() ([]Process, error) {
	items, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var processes []Process
	for _, item := range items {
		if item.IsDir() {
			pid, err := strconv.ParseInt(item.Name(), 10, 64)
			if err == nil {
				proc, err := readProcessInfo(pid)
				if err != nil {
					continue
				}
				processes = append(processes, proc)
			}
		}
	}
	return processes, nil
}

func readProcessInfo(pid int64) (Process, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return Process{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var name string
	var mem float64
	var foundMem bool // Kernel threads don't have VmRSS
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		if strings.HasPrefix(scanner.Text(), "Name:") {
			name = fields[1]
		}
		if strings.HasPrefix(scanner.Text(), "VmRSS:") { // Current Resident Set Size (how much physical RAM the process is currently using).
			parsed, err := strconv.ParseFloat(fields[1], 64)
			if err == nil {
				mem = parsed
				foundMem = true
			}
		}
		if name != "" && foundMem {
			break
		}
	}
	return Process{PID: pid, Name: name, Mem: mem / 1024}, nil
}
