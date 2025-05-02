package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"the-autoscaler/docker"
	pb "the-autoscaler/proto"
	"the-autoscaler/utils"
	"time"

	"google.golang.org/grpc"
)

const (
	maxReplicas    = 3
	minReplicas    = 1
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
	ContainerID string  `json:"containerId,omitempty"`
}

type Node struct {
	Status *SystemStatus
	Full   bool
}

type server struct {
	pb.UnimplementedTheAutocalerServer
}

func createInstance() {
	containerList, err := docker.GetAllContainerIDs()
	if err != nil {
		log.Fatal("Error listing all nodes: ", err)
	}

	if len(containerList) == maxReplicas {
		return
	}

	newNodeName := fmt.Sprintf("node-%v", utils.RandomString(10))

	node, err := docker.CreateInstance(newNodeName)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Node %s created!", node.ID)
	nodeList[node.ID] = &Node{}
}

func deleteInstance(id string) {
	containerList, err := docker.GetAllContainerIDs()
	if err != nil {
		log.Fatal("Error listing all nodes: ", err)
	}

	if len(containerList) == minReplicas {
		return
	}

	if err := docker.DeleteInstance(id); err != nil {
		log.Fatal(err)
	}
	log.Printf("Node %s deleted!", id)
}

var nodeList = map[string]*Node{}

func (s *server) CheckResourcesUsage(_ context.Context, in *pb.CheckRequest) (*pb.CheckReply, error) {
	log.Printf("[%v] CPU: %f Memory: %f\n", in.ContainerId, in.CpuPercent, in.MemoryPercent)

	nodeList[in.ContainerId] = &Node{Status: &SystemStatus{in.CpuPercent, uint64(in.MemoryPercent), in.ContainerId}}

	if in.CpuPercent > CPULimit {
		log.Println("CPU usage is too high.")
		nodeList[in.ContainerId].Full = true
	}

	if in.MemoryPercent < MemoryLimit {
		log.Println("Memory usage is too high.")
		nodeList[in.ContainerId].Full = true
	}

	if in.CpuPercent < CPURequired && in.MemoryPercent > MemoryRequired {
		log.Println("Resources usage is low.")
		nodeList[in.ContainerId].Full = false
		if len(nodeList) > minReplicas {
			deleteInstance(in.ContainerId)
			delete(nodeList, in.ContainerId)
		}
	}

	if nodeList[in.ContainerId].Full && len(nodeList) < maxReplicas {
		createInstance()
	}

	return &pb.CheckReply{}, nil
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterTheAutocalerServer(serv, &server{})

	log.Printf("Starting Orchestrator gRPC Server on :%d!\n", *port)

	if len(nodeList) == 0 {
		log.Println("No nodes in the pool, creating nodes...")
		for len(nodeList) <= minReplicas {
			createInstance()
		}
	}

	err = serv.Serve(list)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
