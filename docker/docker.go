package docker

import (
	"context"
	"fmt"

	client "github.com/influxdata/influxdb1-client"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func CreateInstance(client *client.Client, containerName string) (*container.CreateResponse, error) {
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

func DeleteInstance(client *client.Client, containerID string) error {
	if err := client.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("Error stopping node. %s", err.Error())
	}

	if err := client.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing node. %s", err.Error())
	}

	return nil
}
