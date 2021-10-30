package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
)

const goServices = "go_services"

type Context struct {
	Name     string
	BasePath string
	GoRoot   string // Path to go_services
	Config   config.WorkspaceConfig
}

func InitContext(workspacePath string) (Context, error) {
	if len(workspacePath) == 0 {
		var err error
		workspacePath, err = config.FindWorkspaceConfigPath()
		if err != nil {
			return Context{}, err
		}
	}
	fmt.Printf("workspacePath %s\n", workspacePath)
	fmt.Printf("go root %s\n", filepath.Join(workspacePath, goServices))
	conf, err := config.ReadWorkspaceConfig(workspacePath)
	if err != nil {
		return Context{}, err
	}

	res := Context{
		Name:     filepath.Base(workspacePath), // TODO: validate
		BasePath: workspacePath,
		GoRoot:   filepath.Join(workspacePath, goServices),
		Config:   conf,
	}

	return res, nil
}
