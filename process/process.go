package process

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Process struct {
	PID        int64
	Name       string
	State      string
	CPU        float64
	Mem        float64
	NumThreads int
}

type ProcessTracker struct {
	prevTicks    map[int64]float64 // PID -> previous utime+stime
	prevCPUTotal float64           // prev total CPU from /proc/stat
}

type procResult struct {
	proc     Process
	cpuTicks float64
	pid      int64
}

type SortBy int

const (
	SortByCPU SortBy = iota
	SortByMem
)

var pageSize = os.Getpagesize()

func (pt *ProcessTracker) GetProcessList(curCPUTotal float64, sortBy SortBy) ([]Process, error) {
	deltaTotal := curCPUTotal - pt.prevCPUTotal

	items, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	processes := make([]Process, 0, len(items))
	currentPIDs := make(map[int64]bool) //delete stale pids

	var pids []int64
	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		pid, err := strconv.ParseInt(item.Name(), 10, 64)
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}

	resultCh := make(chan procResult, len(pids))
	var wg sync.WaitGroup

	for _, pid := range pids {
		wg.Add(1)
		go func(pid int64) {
			defer wg.Done()

			proc, cputicks, err := readProcess(pid)
			if err != nil {
				return //process died - skip silently
			}
			resultCh <- procResult{proc: proc, cpuTicks: cputicks, pid: pid}
		}(pid)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {

		if i, ok := pt.prevTicks[result.pid]; ok && deltaTotal > 0 {
			result.proc.CPU = ((result.cpuTicks - i) / (deltaTotal)) * 100
		}

		pt.prevTicks[result.pid] = result.cpuTicks
		currentPIDs[result.pid] = true

		processes = append(processes, result.proc)
	}

	pt.prevCPUTotal = curCPUTotal
	for pid := range pt.prevTicks {
		if !currentPIDs[pid] {
			delete(pt.prevTicks, pid)
		}
	}
	sort.Slice(processes, func(i, j int) bool {
		switch sortBy {
		case SortByMem:
			return processes[i].Mem > processes[j].Mem
		default:
			return processes[i].CPU > processes[j].CPU
		}
	})
	return processes, nil
}

func readProcess(pid int64) (Process, float64, error) {

	b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return Process{}, 0, err
	}
	s := string(b)

	// Process name can contain spaces
	endIndex := strings.LastIndex(s, ")")
	if endIndex == -1 {
		return Process{}, 0, fmt.Errorf("malformed stat for pid %d", pid)
	}
	startIndex := strings.Index(s, "(")
	if startIndex == -1 {
		return Process{}, 0, fmt.Errorf("malformed stat for pid %d", pid)
	}
	name := string(s[startIndex+1 : endIndex])

	rest := s[endIndex+1:]

	fields := strings.Fields(rest)
	if len(fields) < 22 {
		return Process{}, 0, fmt.Errorf("malformed stat for pid %d", pid)
	}

	state := fields[0]

	utime, err := strconv.ParseFloat(fields[11], 64)
	if err != nil {
		return Process{}, 0, err
	}

	stime, err := strconv.ParseFloat(fields[12], 64)
	if err != nil {
		return Process{}, 0, err
	}

	threads, err := strconv.Atoi(fields[17])
	if err != nil {
		return Process{}, 0, err
	}

	rss, err := strconv.ParseFloat(fields[21], 64)
	if err != nil {
		return Process{}, 0, err
	}

	return Process{
		PID:        pid,
		Name:       name,
		State:      state,
		NumThreads: threads,
		Mem:        rss * float64(pageSize) / 1024 / 1024,
	}, utime + stime, nil
}

func NewProcessTracker() *ProcessTracker {
	return &ProcessTracker{
		prevTicks: make(map[int64]float64),
	}
}
