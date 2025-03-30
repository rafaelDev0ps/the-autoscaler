package utils

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func CheckFreeMemory() (uint64, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return mem.Available, nil
}

func CheckCPUPercent() (float64, error) {
	percent, err := cpu.Percent(time.Minute, false)
	if err != nil {
		return 0.0, err
	}

	return percent[0], nil
}
