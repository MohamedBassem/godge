package godge

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

type Executor interface {
	SetDockerClient(*docker.Client)
	Execute(args []string) error
	Stop() error
}

type baseExecutor struct {
	dockerClient *docker.Client
	container    *docker.Container
}

func (b *baseExecutor) SetDockerClient(d *docker.Client) {
	b.dockerClient = d
}

type goExecutor struct {
	baseExecutor
	packgeArchive []byte `json:"packageArchive"`
}

func (g *goExecutor) Execute(args []string) error {
	if g.dockerClient == nil {
		// Panic if there's a logic error
		panic("Docker client must be set for go executor")
	}

	pdir, err := unzipToTmpDir(g.packgeArchive)
	if err != nil {
		return fmt.Errorf("failed to unzip package: %v", err)
	}

	cmd := []string{"go-wrapper", "download", "&&", "go-wrapper", "install", "&&", "go-wrapper", "run"}
	cmd = append(cmd, args...)
	option := docker.CreateContainerOptions{
		Name: randomString(20),
		Config: &docker.Config{
			Image:      "golang:1.8",
			Cmd:        cmd,
			WorkingDir: "/go/src/app",
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				fmt.Sprintf("%v:/go/src/app", pdir),
			},
		},
	}

	g.container, err = g.dockerClient.CreateContainer(option)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}
	if err := g.dockerClient.StartContainer(g.container.ID, nil); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}
	return nil
}

func (g *goExecutor) Stop() error {
	if err := g.dockerClient.StopContainer(g.container.ID, 2); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}
	return nil
}
