package sysinfo

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetUptime() (time.Duration, error) {
	items, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(items))
	t, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(t * float64(time.Second)), nil
}

func GetLoadAvg() ([3]float64, error) {
	items, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return [3]float64{}, err
	}
	s := string(items)

	fields := strings.Fields(s)
	avg1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return [3]float64{}, err
	}
	avg5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return [3]float64{}, err
	}
	avg15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return [3]float64{}, err
	}
	return [3]float64{avg1, avg5, avg15}, nil
}
