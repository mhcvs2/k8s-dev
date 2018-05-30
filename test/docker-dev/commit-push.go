package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	commitOptions := types.ContainerCommitOptions{
		Reference: "new_image:new_tag",
		Comment: "message",
		Pause: true,
		Author: "mhc",
	}

	containerID := "cf5e220dd86c50548e4488c2cbe947be1abd0076ee095990f9a17567621a7636"

	_, err = cli.ContainerCommit(ctx, containerID, commitOptions)

	cli.ImagePush()
}

