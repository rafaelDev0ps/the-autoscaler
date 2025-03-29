package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	maxReplicas = 2
	minReplicas = 1
)

func createInstance(client *client.Client, containerName string) (*container.CreateResponse, error) {
	containerConf := container.Config{
		Image: "server",
	}
	hostConf := container.HostConfig{}

	network := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"bridge": {
				IPAMConfig: &network.EndpointIPAMConfig{},
			},
		},
	}

	platform := v1.Platform{}

	resp, err := client.ContainerCreate(context.Background(), &containerConf, &hostConf, &network, &platform, containerName)
	if err != nil {
		return nil, fmt.Errorf("error creating node. %s", err.Error())
	}

	if err := client.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("error starting node. %s", err.Error())
	}

	return &resp, nil
}

func deleteInstance(client *client.Client, containerID string) error {
	if err := client.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("Error stopping node. %s", err.Error())
	}

	if err := client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing node. %s", err.Error())
	}

	return nil
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
	node, err := createInstance(dockerClient, newNodeName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Node %s created!", node.ID)

	if len(containerList) > maxReplicas {
		for i := 0; i < len(containerList)-maxReplicas; i++ {
			if err := deleteInstance(dockerClient, containerList[i].ID); err != nil {
				log.Fatal("Error deleting replica", err)
			}
		}
	}
	if len(containerList) < minReplicas {
		for i := 0; i < minReplicas-len(containerList); i++ {
			newNode, err := createInstance(dockerClient, newNodeName)
			fmt.Printf("Node %s created!", newNode.ID)
			if err != nil {
				log.Fatal("Error creating replica", err)
			}
		}
	}
}
