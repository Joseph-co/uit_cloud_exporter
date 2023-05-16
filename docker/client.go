package docker

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	dclient "github.com/docker/docker/client"
)

var (
	dockerClient     *dclient.Client
	dockerClientErr  error
	dockerClientOnce sync.Once
)

const dockerEndpoint = "unix:///var/run/docker.sock"
//const dockerMac = "unix:///Users/slc/.docker/run/docker.sock"

var dockerTimeout = 10 * time.Second

func defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), dockerTimeout)
	return ctx
}

// Client creates a Docker API client based on the given Docker flags
func Client() (*dclient.Client, error) {
	dockerClientOnce.Do(func() {
		var client *http.Client
		dockerClient, dockerClientErr = dclient.NewClientWithOpts(
			dclient.WithHost(dockerEndpoint),
			//dclient.WithHost(dockerMac),
			dclient.WithHTTPClient(client),
			dclient.WithAPIVersionNegotiation())
	})
	return dockerClient, dockerClientErr
}
//func GetDockerInfo() types.Info {
//	c, _ := Client()
//	info, err := c.Info(defaultContext())
//	if err != nil {
//		log.Fatal(err)
//	}
//	return info
//}

func GetContainerInspect(ID string) types.ContainerJSON {
	c, _ := Client()
	info, err := c.ContainerInspect(defaultContext(), ID)
	//info, err := c.ContainerInspect(defaultContext(), "9594a9eb913f")
	if err != nil {
		log.Fatal(err)
	}
	return info
}

func GetContainerIDs()[]string {
	c, _ := Client()
	options := types.ContainerListOptions{
		All: true,
	}
	containers, err := c.ContainerList(defaultContext(), options)
	if err != nil {
		log.Fatal(err)
	}
	containerIDs := make([]string, 0, len(containers))
	for _,container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}
	return containerIDs
}
