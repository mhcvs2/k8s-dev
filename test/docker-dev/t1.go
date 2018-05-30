package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"fmt"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Println("------------------------")
		fmt.Println("container id: ", container.ID)
		fmt.Println("image: ", container.Image)
	}
}

