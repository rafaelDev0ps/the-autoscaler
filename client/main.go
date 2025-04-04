package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"the-autoscaler/utils"
	"time"
)

const (
	CPULimit       = 80.0
	MemoryLimit    = 10000000 //10MB
	CPURequired    = 20.0
	MemoryRequired = 1000000 //1MB
	CheckInterval  = 2 * time.Minute
	APIPort        = ":8081"
)

type SystemStatus struct {
	CPUPercent  float64 `json:"cpuPercent"`
	FreeMemory  uint64  `json:"freeMemory"`
	NeedsNode   bool    `json:"needsNode"`
	ContainerID string  `json:"containerId,omitempty"`
}

func requestNewNode() error {
	resp, err := http.Get("http://localhost:8080/create")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to request new node: %s", resp.Status)
	}

	return nil
}

func requestDeleteNode(containerID string) error {
	resp, err := http.Get("http://localhost:8080/delete?id=" + containerID)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to request delete node: %s", resp.Status)
	}

	return nil
}

func getContainerID() (string, error) {
	cmd := exec.Command("cat /containerid")

	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting the container ID. %s", err)
	}

	containerID := out.String()
	return containerID, nil
}

func checkSystem() SystemStatus {
	var status SystemStatus
	var err error

	status.CPUPercent, err = utils.CheckCPUPercent()
	if err != nil {
		log.Printf("Error checking CPU percent: %v", err)
	}

	status.FreeMemory, err = utils.CheckFreeMemory()
	if err != nil {
		log.Printf("Error checking free memory: %v", err)
	}

	if status.CPUPercent > CPULimit {
		log.Println("CPU usage is too high.")
		status.NeedsNode = true
	}

	if status.FreeMemory < MemoryLimit {
		log.Println("Memory usage is too high.")
		status.NeedsNode = true
	}

	if status.CPUPercent < CPURequired && status.FreeMemory > MemoryRequired {
		log.Println("Resources usage is low.")
		status.NeedsNode = false
	}

	if status.NeedsNode {
		err := requestNewNode()
		if err != nil {
			log.Printf("Error requesting new node: %v", err)
		} else {
			log.Println("New node requested.")
		}
	}
	if !status.NeedsNode {
		err := requestDeleteNode(status.ContainerID)
		if err != nil {
			log.Printf("Error requesting delete node: %v", err)
		} else {
			log.Println("Node deleted.")
		}
	}

	status.ContainerID, err = getContainerID()
	if err != nil {
		log.Printf("Error getting container ID: %v", err)
	} else {
		log.Println("Container ID:", status.ContainerID)
	}

	return status
}

func main() {
	// Set up periodic checks
	go func() {
		ticker := time.NewTicker(CheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Running scheduled system check...")
				checkSystem()
			}
		}
	}()

	log.Printf("Starting API server on port %s...", APIPort)
	log.Printf("System checks will run every %v", CheckInterval)

	if err := http.ListenAndServe(APIPort, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
