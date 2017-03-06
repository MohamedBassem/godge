package godge

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	docker "github.com/fsouza/go-dockerclient"
)

type Executor interface {
	setDockerClient(*docker.Client)
	Execute(args []string) error
	ReadFileFromContainer(path string) (string, error)
	Stdout() (string, error)
	Stderr() (string, error)
	Stop() error
}

type baseExecutor struct {
	dockerClient *docker.Client
	container    *docker.Container
	workDir      string
	stoppedOnce  sync.Once
}

func (b *baseExecutor) setDockerClient(d *docker.Client) {
	b.dockerClient = d
}

func (b *baseExecutor) ReadFileFromContainer(path string) (string, error) {
	buf := new(bytes.Buffer)
	option := docker.DownloadFromContainerOptions{
		OutputStream: buf,
		Path:         fmt.Sprintf("%v/%v", b.workDir, path),
	}
	if err := b.dockerClient.DownloadFromContainer(b.container.ID, option); err != nil {
		return "", fmt.Errorf("failed to read file from container: %v", err)
	}
	return string(buf.Bytes()), nil
}

func (b *baseExecutor) Stdout() (string, error) {
	buf := new(bytes.Buffer)
	option := docker.LogsOptions{
		OutputStream: buf,
		Container:    b.container.ID,
		Stdout:       true,
		Tail:         "all",
	}
	if err := b.dockerClient.Logs(option); err != nil {
		return "", fmt.Errorf("failed to read stdout from container: %v", err)
	}
	return string(buf.Bytes()), nil
}

func (b *baseExecutor) Stderr() (string, error) {
	buf := new(bytes.Buffer)
	option := docker.LogsOptions{
		ErrorStream: buf,
		Container:   b.container.ID,
		Stderr:      true,
		Tail:        "all",
	}
	if err := b.dockerClient.Logs(option); err != nil {
		return "", fmt.Errorf("failed to read stderr from container: %v", err)
	}
	return string(buf.Bytes()), nil
}

func (g *baseExecutor) Stop() error {
	var err error
	g.stoppedOnce.Do(func() {
		if err = g.dockerClient.StopContainer(g.container.ID, 2); err != nil {
			err = fmt.Errorf("failed to stop container: %v", err)
			return
		}
	})
	return err
}

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
