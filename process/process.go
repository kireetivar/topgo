package process

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Process struct {
	PID  int64
	Name string
	CPU  float64
	Mem  float64
}

type ProcessTracker struct {
	prevTicks    map[int64]float64 // PID -> previous utime+stime
	prevCPUTotal float64           // prev total CPU from /proc/stat
}

func (pt *ProcessTracker) GetProcessList(curCPUTotal float64) ([]Process, error) {
	deltaTotal := curCPUTotal - pt.prevCPUTotal

	items, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var processes []Process
	currentPIDs := make(map[int64]bool) //delete stale pids

	for _, item := range items {
		if item.IsDir() {
			pid, err := strconv.ParseInt(item.Name(), 10, 64)
			if err == nil {
				proc, err := readProcessInfo(pid)
				if err != nil {
					continue
				}
				cpuTicks, err := readProcessCPUTicks(pid)
				if err != nil {
					continue
				}

				if i, ok := pt.prevTicks[pid]; ok && deltaTotal > 0 {
					proc.CPU = ((cpuTicks - i) / (deltaTotal)) * 100
				}

				pt.prevTicks[pid] = cpuTicks
				currentPIDs[pid] = true

				processes = append(processes, proc)
			}
		}
	}
	pt.prevCPUTotal = curCPUTotal
	for pid := range pt.prevTicks {
		if !currentPIDs[pid] {
			delete(pt.prevTicks, pid)
		}
	}
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU
	})
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

func readProcessCPUTicks(pid int64) (float64, error) {
	b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	s := string(b)

	// Process name can contain spaces
	index := strings.LastIndex(s, ")")
	if index == -1 {
		return 0, fmt.Errorf("malformed stat for pid %d", pid)
	}

	rest := s[index+1:]

	fields := strings.Fields(rest)
	if len(fields) < 13 {
		return 0, fmt.Errorf("malformed stat for pid %d", pid)
	}

	utime, err := strconv.ParseFloat(fields[11], 64)
	if err != nil {
		return 0, err
	}

	stime, err := strconv.ParseFloat(fields[12], 64)
	if err != nil {
		return 0, err
	}

	return utime + stime, nil
}

func NewProcessTracker() *ProcessTracker {
	return &ProcessTracker{
		prevTicks: make(map[int64]float64),
	}
}
