package godge

type Executor interface {
	Execute(args []string) error
	Stop() error
}

type goExecutor struct {
	Package string `json:"package"`
}

func (g *goExecutor) Execute(args []string) error {
	return nil
}

func (g *goExecutor) Stop() error {
	return nil
}
