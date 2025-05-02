package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	pb "the-autoscaler/proto"
	"the-autoscaler/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SystemStatus struct {
	CPUPercent  float64 `json:"cpuPercent"`
	FreeMemory  uint64  `json:"freeMemory"`
	ContainerID string  `json:"containerId,omitempty"`
}

func getContainerID() (string, error) {
	cmd := exec.Command("cat", "/containerid")

	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error getting the container ID. %s", err)
	}

	containerID := out.String()
	return containerID, nil
}

func checkSystem() (*SystemStatus, error) {
	status := &SystemStatus{}
	var err error

	status.ContainerID, err = getContainerID()
	if err != nil {
		return nil, fmt.Errorf("error getting container ID: %v", err)
	} else {
		log.Println("Container ID:", status.ContainerID)
	}

	status.CPUPercent, err = utils.CheckCPUPercent()
	if err != nil {
		return nil, fmt.Errorf("error checking CPU percent: %v", err)
	}

	status.FreeMemory, err = utils.CheckFreeMemory()
	if err != nil {
		return nil, fmt.Errorf("error checking free memory: %v", err)
	}

	return status, nil
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()

	log.Printf("Starting gRPC client on port %s...\n", *addr)
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewTheAutocalerClient(conn)

	// collect metrics every minute
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			systemStatus, err := checkSystem()
			if err != nil {
				log.Fatalf("failed checking client system status: %v", err)
			}

			request := &pb.CheckRequest{
				CpuPercent:    systemStatus.CPUPercent,
				MemoryPercent: float64(systemStatus.FreeMemory),
				ContainerId:   systemStatus.ContainerID,
			}

			_, err = client.CheckResourcesUsage(ctx, request)
			if err != nil {
				log.Fatalf("failed getting client resources usage: %v", err)
			}
		}()
	}
}
