package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"the-autoscaler/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	maxReplicas = 2
	minReplicas = 1
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error with docker client: ", err)
	}

	defer dockerClient.Close()

	containerList, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatal("Error listing all nodes: ", err)
	}

	newNodeName := fmt.Sprintf("node%v", len(containerList)+1)
	node, err := docker.CreateInstance(dockerClient, newNodeName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Node %s created!", node.ID)

	if len(containerList) > maxReplicas {
		for i := 0; i < len(containerList)-maxReplicas; i++ {
			if err := docker.DeleteInstance(dockerClient, containerList[i].ID); err != nil {
				log.Fatal("Error deleting replica", err)
			}
		}
	}
	if len(containerList) < minReplicas {
		for i := 0; i < minReplicas-len(containerList); i++ {
			newNode, err := docker.CreateInstance(dockerClient, newNodeName)
			log.Printf("Node %s created!", newNode.ID)

			if err != nil {
				log.Fatal("Error creating replica", err)
			}
		}
	}

	http.HandleFunc("/health", healthCheckHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
