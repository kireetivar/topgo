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

func (c *CPUStat) GetCPUUsage(curtotal float64, curidle float64) (float64, error) {
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

func ReadTotalCPUTicks() (float64, float64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var curTotal float64
	var curIdle float64

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "cpu" {
			continue
		}

		for i := 1; i < len(fields); i++ {
			val, err := strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse cpu field: %w", err)
			}
			if i == 4 || i == 5 {
				curIdle += val // idle and iowait
			}
			curTotal += val
		}

		break
	}
	return curTotal, curIdle, nil
}
