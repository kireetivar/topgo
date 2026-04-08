package process

import (
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
				if err == nil {
					processes = append(processes, proc)
				}
			}
		}
	}
	return processes, nil
}

func readProcessInfo(pid int64) (Process, error) {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return Process{}, err
	}
	fields := strings.Fields(string(stat))
	var name string
	var mem float64
	for i, field := range fields {
		if i == 1 {
			name = strings.Trim(field, "()")
		}
		if i == 23 {
			parsed, err := strconv.ParseFloat(field, 64)
			if err == nil {
				mem = parsed
			}
		}
	}
	return Process{PID: pid, Name: name, Mem: mem / 1024}, nil
}