package godge

import (
	"bytes"
	"fmt"
	"sync"

	docker "github.com/fsouza/go-dockerclient"
)

// Executor is used to interact with the submission.
type Executor interface {
	setDockerClient(*docker.Client)
	// Excutes the submitted code with the provided arguments.
	Execute(args []string) error
	// Reads a certain file from the container's workspace.
	ReadFileFromContainer(path string) (string, error)
	// Returns the contents of the stdout of the container.
	Stdout() (string, error)
	// Returns the contents for the stderr of the container.
	Stderr() (string, error)
	// Stops the running binary.
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

// ReadFileFromContainer reads a certain file from the container's workspace. The path
// is relative to the container's workdir.
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

// Stdout returns the content of the stdout of the container.
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

// Stderr returns the content of the stderr of the container.
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

// Stop stops the running binary.
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
