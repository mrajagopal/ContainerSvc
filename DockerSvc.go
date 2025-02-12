package containersvc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerService struct {
	cli *client.Client
}

// type Neo4jGraphRepo struct {
// 	drv neo4j.DriverWithContext
// }

func NewDockerService() (*DockerService, error) {
	// Create a new Docker client
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}
	return &DockerService{cli: client}, nil
}

func (d *DockerService) CreateContainer() error {
	// // Create a new Docker client
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	return fmt.Errorf("failed to create Docker client: %v", err)
	// }

	ctx := context.Background()

	// Define the ClickHouse container options
	containerName := "xray-clickhouse-server"
	imageName := "clickhouse/clickhouse-server:latest"
	containerConfig := &container.Config{
		Image: imageName,
		Cmd:   []string{"/entrypoint.sh"}, // Specify the default entrypoint
		ExposedPorts: nat.PortSet{
			"9000/tcp": {},
			"8123/tcp": {},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"9000/tcp": {{HostPort: "9001"}},
			"8123/tcp": {{HostPort: "8124"}},
		},
	}

	// Pull options with authentication
	pullOpts := types.ImagePullOptions{
		All: true,
	}

	_, err := d.cli.ImagePull(ctx, imageName, pullOpts)
	if err != nil {
		return fmt.Errorf("failed to pull ClickHouse image: %v", err)
	}

	// Wait a moment to ensure the pull has completed
	time.Sleep(2 * time.Second)

	// Create the container
	log.Println("Creating ClickHouse container...")
	resp, err := d.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("failed to create ClickHouse container: %v", err)
	}

	// Start the container
	log.Println("Starting ClickHouse container...")
	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start ClickHouse container: %v", err)
	}

	log.Printf("ClickHouse server is now running in container '%s'.\n", containerName)
	return nil
}
