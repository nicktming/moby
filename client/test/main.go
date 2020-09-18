package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	dclient "github.com/docker/docker/client"
	"sync"
	"github.com/docker/go-connections/tlsconfig"
	"net/http"
)

const DefaultDockerHost = "unix:///var/run/docker.sock"

func main() {
	//cli, err := dclient.NewClientWithOpts(dclient.FromEnv)
	cli, err := ClientWithAuth()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

var (
	dockerClient     *dclient.Client
	dockerClientErr  error
	dockerClientOnce sync.Once
)

// Client creates a Docker API client based on the given Docker flags
func ClientWithAuth() (*dclient.Client, error) {
	dockerClientOnce.Do(func() {
		var client *http.Client
		client = &http.Client{}
		options := tlsconfig.Options{
			CAFile:             "/etc/kubernetes/pki/ca.pem",
			CertFile:           "/etc/kubernetes/pki/docker.pem",
			KeyFile:            "/etc/kubernetes/pki/docker-key.pem",
			InsecureSkipVerify: false,
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			dockerClientErr = err
			return
		}
		client.Transport = &http.Transport{
			TLSClientConfig: tlsc,
		}

		dockerClient, dockerClientErr = dclient.NewClient(DefaultDockerHost,
			"",
			client,
			nil)

	})
	return dockerClient, dockerClientErr
}


