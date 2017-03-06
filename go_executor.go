package godge

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type GoExecutor struct {
	baseExecutor
	PackageArchive []byte `json:"packageArchive"`
}

func (g *GoExecutor) Execute(args []string) error {
	if g.dockerClient == nil {
		// Panic if there's a logic error
		panic("Docker client must be set for go executor")
	}

	pdir, err := unzipToTmpDir(g.PackageArchive)
	if err != nil {
		return fmt.Errorf("failed to unzip package: %v", err)
	}

	cmd := []string{"/bin/bash", "-c", fmt.Sprintf(`
	set -e;
	go-wrapper download > /dev/null 2>&1;
	go-wrapper install > /dev/null 2>&1;
	app %v;`, strings.Join(args, " "))}
	cmd = append(cmd, args...)
	wdir := "/go/src/app"
	g.workDir = wdir
	option := docker.CreateContainerOptions{
		Name: randomString(20),
		Config: &docker.Config{
			Image:      "golang:1.8",
			Cmd:        cmd,
			WorkingDir: wdir,
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				fmt.Sprintf("%v:%v", pdir, wdir),
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
