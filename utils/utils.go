package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func readCPUTimes() []int64 {
	// TODO: read /proc/stat instead dummy_stat
	statFile, err := os.Open("dummy_stat")
	if err != nil {
		log.Fatal(err)
	}

	var cpuTimes []string
	scanner := bufio.NewScanner(statFile)
	for scanner.Scan() {
		line := scanner.Text()
		props := strings.Fields(line)
		if props[0] == "cpu" {
			cpuTimes = props[1:]
		}
	}

	cpuTimesFormatted := []int64{}
	for _, attr := range cpuTimes {
		attrFmt, err := strconv.ParseInt(attr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		cpuTimesFormatted = append(cpuTimesFormatted, attrFmt)
	}

	return cpuTimesFormatted
}

func checkCPUPercent(timeInterval int) float64 {
	cpuTime1 := readCPUTimes()
	time.Sleep(time.Duration(timeInterval) * time.Second)
	cpuTime2 := readCPUTimes()

	delta := [10]int64{}
	for i := 0; i < len(cpuTime1)-1; i++ {
		delta[i] = cpuTime2[i] - cpuTime1[i]
	}

	idle := delta[3]
	iowait := delta[4]

	var totalTime int64 = 0
	for i := 0; i < len(delta); i++ {
		totalTime = totalTime + delta[i]
	}

	if totalTime == 0 {
		return 0.0
	}

	totalIdle := idle - iowait

	totalActive := totalTime - totalIdle

	cpuPercent := float64(totalActive) / float64(totalTime) * 100

	return cpuPercent
}

func checkFreeMemory() int64 {

	memInfoFile, err := os.Open("dummy_meminfo")
	if err != nil {
		log.Fatal(err)
	}

	var memFree string
	scanner := bufio.NewScanner(memInfoFile)
	for scanner.Scan() {
		line := scanner.Text()
		props := strings.Fields(line)

		if props[0] == "MemFree:" {
			memFree = props[1]
			break
		}
	}

	defer memInfoFile.Close()

	free, err := strconv.ParseInt(memFree, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return free
}

func main() {
	free := checkFreeMemory()
	log.Println("Free memory: ", free)

	percent := checkCPUPercent(2)
	log.Println("CPU Usage: ", percent)
}
