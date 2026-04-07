package process

import (
	"os"
	"strconv"
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
		return []Process{}, err
	}
	var processes []Process
	for _, item := range items {
		if item.IsDir() {
			pid, err := strconv.ParseInt(item.Name(), 10, 64)
			if err == nil {
				processes = append(processes, Process{PID: pid})
			}
		}
	}
	return processes, nil
}