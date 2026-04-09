package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUStat struct {
	prevIdle  float64
	prevTotal float64
}

func (c *CPUStat) GetCPUUsage() (float64, error) {
	osFile, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer osFile.Close()

	scanner := bufio.NewScanner(osFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "cpu" {
			continue
		}

		// cpu  88 0 140 42651 83 0 17 0 0 0
		cpuNums := make([]float64, len(fields)-1)
		for i := 1; i < len(fields); i++ {
			val, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse cpu field: %w", err)
			}
			cpuNums[i-1] = val
		}

		curidle := cpuNums[3] + cpuNums[4]

		curtotal := 0.0
		for _, num := range cpuNums {
			curtotal += num
		}

		if c.prevTotal == 0 {
			c.prevIdle = curidle
			c.prevTotal = curtotal
			return 0, nil
		}

		deltaIdle := curidle - c.prevIdle
		deltaTotal := curtotal - c.prevTotal

		c.prevIdle = curidle
		c.prevTotal = curtotal

		if deltaTotal == 0 {
			return 0, nil
		}
		return ((deltaTotal - deltaIdle) / deltaTotal) * 100, nil
	}
	return 0, nil
}
