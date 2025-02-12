package containersvc

// This program creates and starts a ClickHouse server in a Docker container using the Docker API.
type ContainerService interface {
	CreateContainer() error
}
