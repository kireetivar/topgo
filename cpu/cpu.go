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

func (c *CPUStat) GetCPUUsage() float64 {
	osFile, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Printf("ERROR: failed to open /proc/stat: %v\n", err)
		return 0
	}
	defer osFile.Close()

	scanner := bufio.NewScanner(osFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if fields[0] != "cpu" {
			continue
		}

		// cpu  88 0 140 42651 83 0 17 0 0 0
		cpuNums := make([]float64, len(fields)-1)
		for i := 1; i < len(fields); i++ {
			cpuNums[i-1] = strToFloat64(fields[i])
		}

		curidle := cpuNums[3] + cpuNums[4]

		curtotal := 0.0
		for _, num := range cpuNums {
			curtotal += num
		}

		if c.prevTotal == 0 {
			c.prevIdle = curidle
			c.prevTotal = curtotal
			return 0
		}

		deltaIdle := curidle - c.prevIdle
		deltaTotal := curtotal - c.prevTotal

		c.prevIdle = curidle
		c.prevTotal = curtotal

		if deltaTotal == 0 {
			return 0
		}
		return ((deltaTotal - deltaIdle) / deltaTotal) * 100
	}
	return 0
}

func strToFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}
