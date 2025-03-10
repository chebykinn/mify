package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/service/tpl"
)

//go:embed tpl/api.yaml.tpl
var apiSchemaTemplate string

func CreateService(mutContext *mutators.MutatorContext, language mifyconfig.ServiceLanguage, serviceName string) error {
	mutContext.GetLogger().Printf("Creating service '%s' ...", serviceName)

	return createServiceImpl(mutContext, language, serviceName, true)
}

func CreateFrontend(mutContext *mutators.MutatorContext, template string, name string) error {
	mutContext.GetLogger().Printf("Creating frontend '%s' ...", name)

	if template == "vue_js" {
		return createServiceImpl(mutContext, mifyconfig.ServiceLanguageJs, name, false)
	}

	return fmt.Errorf("unknown template %s", template)
}

func CreateApiGateway(mutContext *mutators.MutatorContext) error {
	exists, err := checkServiceExists(mutContext, workspace.ApiGatewayName)
	if err != nil {
		return fmt.Errorf("can't check if service exists: %w", err)
	}

	if exists {
		return fmt.Errorf("api gateway already exists, skipping creation")
	}

	err = CreateService(mutContext, mifyconfig.ServiceLanguageGo, workspace.ApiGatewayName)
	if err != nil {
		return err
	}

	return nil
}

func createServiceImpl(
	mutContext *mutators.MutatorContext,
	language mifyconfig.ServiceLanguage,
	serviceName string,
	addOpenApi bool) error {

	conf := mifyconfig.ServiceConfig{
		ServiceName: serviceName,
		Language:    language,
	}

	err := conf.Dump(mutContext.GetDescription().GetMifySchemaAbsPath(serviceName))
	if err != nil {
		return err
	}

	if addOpenApi {
		openapiSchemaPath := mutContext.GetDescription().GetApiSchemaAbsPath(serviceName, workspace.MainApiSchemaName)
		err := render.RenderTemplate(
			apiSchemaTemplate,
			tpl.NewApiSchemaModel(serviceName),
			openapiSchemaPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkServiceExists(mutContext *mutators.MutatorContext, serviceName string) (bool, error) {
	schemasDirAbsPath := mutContext.GetDescription().GetApiSchemaDirAbsPath(serviceName)
	if _, err := os.Stat(schemasDirAbsPath); os.IsNotExist(err) {
		return false, nil
	}

	files, err := os.ReadDir(schemasDirAbsPath)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}
