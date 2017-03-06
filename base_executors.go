package godge

import (
	"bytes"
	"fmt"
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

func (b *baseExecutor) Stop() error {
	var err error
	b.stoppedOnce.Do(func() {
		if err = b.dockerClient.StopContainer(b.container.ID, 2); err != nil {
			err = fmt.Errorf("failed to stop container: %v", err)
			return
		}
	})
	return err
}
