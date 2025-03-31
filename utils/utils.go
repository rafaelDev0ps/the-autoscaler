package utils

import (
	"math/rand"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

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

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}
