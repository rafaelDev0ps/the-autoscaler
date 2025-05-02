package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"gopkg.in/pipe.v2"
)

func GetHostIP() (string, error) {
	p := pipe.Line(
		pipe.Exec("ifconfig"),
		pipe.Exec("grep", "-w", "inet"),
		pipe.Exec("awk", "{print $2}"),
		pipe.Exec("tail", "-n", "1"),
	)
	output, err := pipe.CombinedOutput(p)
	if err != nil {
		return "", fmt.Errorf("error getting the host ip. %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func CreateInstance(containerName string) (*container.CreateResponse, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Error with docker client: ", err)
	}

	defer dockerClient.Close()

	hostIP, err := GetHostIP()
	if err != nil {
		return nil, fmt.Errorf("error getting host ip during pod creation: %v", err)
	}

	containerConf := container.Config{
		Image: "client",
		Cmd:   []string{"--addr", hostIP + ":50051"},
	}

	hostConf := container.HostConfig{
		NetworkMode: container.NetworkMode("host"),
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/etc/app/containerid",
				Target: "/containerid",
			},
		},
	}

	// network := network.NetworkingConfig{
	// 	EndpointsConfig: map[string]*network.EndpointSettings{
	// 		"host": {
	// 			IPAMConfig: &network.EndpointIPAMConfig{},
	// 		},
	// 	},
	// }
	network := network.NetworkingConfig{}

	platform := v1.Platform{}

	resp, err := dockerClient.ContainerCreate(context.Background(), &containerConf, &hostConf, &network, &platform, containerName)
	if err != nil {
		return nil, fmt.Errorf("error creating node. %s", err.Error())
	}

	err = os.WriteFile("/etc/app/containerid", []byte(resp.ID), 0644)
	if err != nil {
		log.Fatal("Error writing container ID to file: ", err)
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
