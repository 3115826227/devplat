package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"testing"
)

func TestNewDockerManager(t *testing.T) {
	//fmt.Println(dm.ClientVersion())
	InitDockerManager()
	images, err := dm.ListImages()
	if err != nil {
		return
	}
	imageMap := make(map[string]types.ImageSummary)
	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			imageMap[repoTag] = image
		}
	}
	fmt.Println(imageMap)
}

func TestGetPath(t *testing.T) {
}
