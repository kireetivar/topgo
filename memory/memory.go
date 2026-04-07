package memory

import (
	"os"
	"bufio"
	"strings"
	"strconv"
)

func GetMemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var memTotal, memAvailable int
	var memUsage float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			memTotal, err = strconv.Atoi(parts[1])
			if err != nil {
				return 0, err
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			memAvailable, err = strconv.Atoi(parts[1])
			if err != nil {
				return 0, err
			}
		}
		if memTotal > 0 && memAvailable > 0 {
			memUsage = float64(memTotal - memAvailable)
			break
		}
	}
	return (memUsage / float64(memTotal)) * 100, nil
}
