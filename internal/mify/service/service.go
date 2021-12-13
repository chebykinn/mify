package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/lang"
	"github.com/chebykinn/mify/internal/mify/workspace"
)

const (
	ApiGatewayName string = "api_gateway"
)

func checkServiceExists(basePath string, name string) (bool, error) {
	services, err := filepath.Glob(filepath.Join(basePath, "*_services/*"))
	if err != nil {
		return false, err
	}

	for _, f := range services {
		stat, err := os.Stat(f)
		if err != nil {
			return false, err
		}
		if !stat.IsDir() {
			continue
		}
		svc := filepath.Base(f)
		if svc == name {
			return true, nil
		}
	}

	return false, nil
}

func makeServiceList(basePath string, language lang.ServiceLanguage, name string) ([]string, error) {
	lst := []string{}
	if language == lang.ServiceLanguageJs {
		services, err := filepath.Glob(filepath.Join(basePath, "js_services/*"))
		if err != nil {
			return nil, err
		}
		isAdded := false
		for _, f := range services {
			stat, err := os.Stat(f)
			if err != nil {
				return nil, err
			}
			if !stat.IsDir() {
				continue
			}
			svc := filepath.Base(f)
			if svc == name {
				isAdded = true
			}
			lst = append(lst, svc)
		}
		if !isAdded {
			lst = append(lst, name)
		}
	}

	return lst, nil
}

func CreateService(ctx *core.Context, wspContext workspace.Context, language lang.ServiceLanguage, name string) error {
	fmt.Printf("Creating service: %s\n", name)

	repo := fmt.Sprintf("%s/%s/%s",
		wspContext.Config.GitHost,
		wspContext.Config.GitNamespace,
		wspContext.Config.GitRepository)

	svcList, err := makeServiceList(wspContext.BasePath, language, name)
	if err != nil {
		return err
	}
	context := Context{
		ServiceName: name,
		Repository:  repo,
		Language:    language,
		// TODO: separate lang specific params
		GoModule:    repo + "/go_services",
		Workspace:   wspContext,
		ServiceList: svcList,
	}

	if err := RenderTemplateTree(ctx, context); err != nil {
		return err
	}

	conf := config.ServiceConfig{
		ServiceName: name,
		Language:    language,
	}

	if err := config.SaveServiceConfig(wspContext.BasePath, name, conf); err != nil {
		return err
	}

	return nil
}

func CreateFrontend(ctx *core.Context, wspContext workspace.Context, template string, name string) error {
	if template == "vue_js" {
		CreateService(ctx, wspContext, lang.ServiceLanguageJs, name)
		return nil
	}

	return fmt.Errorf("unknown template %s", template)
}

func TryCreateApiGateway(ctx *core.Context, wspContext workspace.Context) (bool, error) {
	apiGatewayExists, err := checkServiceExists(wspContext.BasePath, ApiGatewayName)
	if err != nil {
		return false, err
	}

	if !apiGatewayExists {
		if err := CreateService(ctx, wspContext, lang.ServiceLanguageGo, ApiGatewayName); err != nil {
			return true, err
		}
	} else {
		fmt.Printf("Api gateway already exists. Skipping creation... \n")
	}

	return false, nil
}
