package godge

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

// GoExecutor implements the Executor interface. It's used in the submit request
// when the language is Go. You won't deal with the GoExecutor directly, it's only
// exposed to be used by the command line client.
type GoExecutor struct {
	baseExecutor
	// A zip archive containing the "main" package to be executed.
	PackageArchive []byte `json:"packageArchive"`
}

// Execute executes the Go main package submitted with the given arguments.
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
