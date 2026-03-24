package mem

import "fmt"
import "os"
import "bufio"
import "strings"
import "strconv"

func GetMemoryUsage() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("Error while opening file: ", err)
		return 0
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
				fmt.Println("Error while converting str to int: ", err)
				return 0
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			memAvailable, err = strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error while conerting str to int: ", err)
				return 0
			}
		}
		if memTotal > 0 && memAvailable > 0 {
			memUsage = float64(memTotal - memAvailable)
			break
		}
	}
	return (memUsage / float64(memTotal)) * 100
}
