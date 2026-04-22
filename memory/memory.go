package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MemStats struct {
	UsagePercentage float64
	TotalGB         float64
	SwapPercentage  float64
	SwapTotalGB     float64
}

func GetMemoryUsage() (MemStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemStats{}, err
	}
	defer file.Close()
	var memTotal, memAvailable, swapTotal, swapFree int
	var memUsage, swapUsage, swapPercentage float64
	found := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			memTotal, err = strconv.Atoi(parts[1])
			if err != nil {
				return MemStats{}, err
			}
			found++
			continue
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			memAvailable, err = strconv.Atoi(parts[1])
			if err != nil {
				return MemStats{}, err
			}
			found++
			continue
		}
		if strings.HasPrefix(line, "SwapTotal:") {
			parts := strings.Fields(line)
			swapTotal, err = strconv.Atoi(parts[1])
			if err != nil {
				return MemStats{}, err
			}
			found++
			continue
		}
		if strings.HasPrefix(line, "SwapFree:") {
			parts := strings.Fields(line)
			swapFree, err = strconv.Atoi(parts[1])
			if err != nil {
				return MemStats{}, err
			}
			found++
			continue
		}
		if found == 4 {
			break
		}
	}
	memUsage = float64(memTotal - memAvailable)
	swapUsage = float64(swapTotal - swapFree)
	if swapTotal > 0 {
		swapPercentage = (swapUsage / float64(swapTotal)) * 100
	}
	return MemStats{
		UsagePercentage: (memUsage / float64(memTotal)) * 100,
		TotalGB:         float64(memTotal) / (1024 * 1024),
		SwapPercentage:  swapPercentage,
		SwapTotalGB:     float64(swapTotal) / (1024 * 1024),
	}, nil
}
