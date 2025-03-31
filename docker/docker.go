package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func CreateInstance(containerName string) (*container.CreateResponse, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error with docker client: ", err)
	}

	defer dockerClient.Close()

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

	resp, err := dockerClient.ContainerCreate(context.Background(), &containerConf, &hostConf, &network, &platform, containerName)
	if err != nil {
		return nil, fmt.Errorf("error creating node. %s", err.Error())
	}

	if err := dockerClient.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("error starting node. %s", err.Error())
	}

	return &resp, nil
}

func DeleteInstance(containerID string) error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error with docker client: ", err)
	}

	defer dockerClient.Close()

	if err := dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("Error stopping node. %s", err.Error())
	}

	if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing node. %s", err.Error())
	}

	return nil
}

func GetAllContainerIDs() ([]string, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error with docker client: ", err)
	}

	defer dockerClient.Close()

	containerList, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("Error listing all nodes: %s", err)
	}

	containerIDs := make([]string, len(containerList))
	for i, container := range containerList {
		containerIDs[i] = container.ID
	}

	return containerIDs, nil
}
