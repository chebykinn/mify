// TODO: refactor!!!

package workspace

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mify-io/mify/internal/mify/util"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

const (
	ApiGatewayName    = "api-gateway"
	MainApiSchemaName = "api.yaml"
	MifySchemaName    = "service.mify.yaml"
	CloudSchemaName   = "cloud.mify.yaml"
	GoServicesDirName = "go-services"
	DevRunnerName     = "dev-runner"
	TmpSubdir         = "services"
)

var (
	ErrUnsupportedLanguage = errors.New("unknown or unsupported language")
	ErrNoSuchService       = errors.New("no such service")
)

type GoService struct {
	Name string
}

type Description struct {
	Name        string
	BasePath    string
	GoRoot      string // Path to go-services
	Config      mifyconfig.WorkspaceConfig
	TplHeader   string
	TplHeaderPy string
}

func InitDescription(workspacePath string) (Description, error) {
	wrapError := func(err error) error {
		return fmt.Errorf("can't initialize description: %w", err)
	}

	if len(workspacePath) == 0 {
		var err error
		workspacePath, err = mifyconfig.FindWorkspaceConfigPath()
		if err != nil {
			return Description{}, wrapError(err)
		}
	}

	conf, err := mifyconfig.ReadWorkspaceConfig(workspacePath)
	if err != nil {
		return Description{}, wrapError(err)
	}

	res := Description{
		Name:        filepath.Base(workspacePath), // TODO: validate
		BasePath:    workspacePath,
		GoRoot:      filepath.Join(workspacePath, mifyconfig.GoServicesRoot),
		Config:      conf,
		TplHeader:   "// THIS FILE IS AUTOGENERATED, DO NOT EDIT\n// Generated by mify",
		TplHeaderPy: "# THIS FILE IS AUTOGENERATED, DO NOT EDIT\n# Generated by mify",
	}

	return res, nil
}

func (c Description) GetApiServices() []string {
	services := []string{}
	files, err := os.ReadDir(c.GetSchemasRootAbsPath())
	if err != nil {
		return nil
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		services = append(services, f.Name())
	}
	return services
}

func (c Description) GetFrontendServices() ([]string, error) {
	services := []string{}
	files, err := os.ReadDir(c.GetSchemasRootAbsPath())
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		cfgPath := c.GetMifySchemaAbsPath(f.Name())
		cfg, err := mifyconfig.ReadServiceCfg(cfgPath)
		if err != nil {
			return nil, err
		}

		if cfg.Language == mifyconfig.ServiceLanguageJs {
			services = append(services, f.Name())
		}
	}
	return services, nil
}

func (c Description) GetAllApps() []string {
	services := c.GetApiServices()
	return append(services, DevRunnerName)
}

func (c Description) HasService(serviceName string) bool {
	services := c.GetApiServices()
	for _, svc := range services {
		if serviceName == svc {
			return true
		}
	}
	return false
}

// Path to include 'app' package
func (c Description) GetAppIncludePath(serviceName string) string {
	return fmt.Sprintf(
		"%s/internal/%s/generated/app",
		c.GetGoModule(),
		serviceName)
}

// Path to include 'core' package
func (c *Description) GetCoreIncludePath(serviceName string) string {
	return fmt.Sprintf(
		"%s/internal/%s/generated/core",
		c.GetGoModule(),
		serviceName)
}

func (c Description) GetSchemasRootRelPath() string {
	return "schemas"
}

func (c Description) GetSchemasRootAbsPath() string {
	return path.Join(c.BasePath, c.GetSchemasRootRelPath())
}

func (c Description) GetSchemasRelPath(serviceName string) string {
	return path.Join(c.GetSchemasRootRelPath(), serviceName)
}

func (c Description) GetSchemasAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetSchemasRelPath(serviceName))
}

func (c Description) GetMifySchemaRelPath(serviceName string) string {
	return path.Join(c.GetSchemasRelPath(serviceName), MifySchemaName)
}

func (c Description) GetMifySchemaAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetMifySchemaRelPath(serviceName))
}

func (c Description) GetCloudSchemaRelPath(serviceName string) string {
	return path.Join(c.GetSchemasRelPath(serviceName), CloudSchemaName)
}

func (c Description) GetCloudSchemaAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetCloudSchemaRelPath(serviceName))
}

func (c Description) GetApiSchemaDirRelPath(serviceName string) string {
	return path.Join("schemas", serviceName, "api")
}

func (c Description) GetApiSchemaDirAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetApiSchemaDirRelPath(serviceName))
}

func (c Description) GetApiSchemaAbsPath(serviceName string, schemaName string) string {
	return path.Join(c.BasePath, "schemas", serviceName, "api", schemaName)
}

// Abs path to api_generated.yaml
func (c Description) GetApiSchemaGenAbsPath(serviceName string) string {
	return path.Join(c.BasePath, "schemas", serviceName, "api/api_generated.yaml")
}

func (c *Description) GetRepository() string {
	return fmt.Sprintf("%s/%s/%s",
		c.Config.GitHost,
		c.Config.GitNamespace,
		c.Config.GitRepository)
}

func (c Description) GetGoModule() string {
	return fmt.Sprintf("%s/%s",
		c.GetRepository(),
		mifyconfig.GoServicesRoot)
}

func (c Description) GetGoConfigsImportPath() string {
	return fmt.Sprintf("%s/%s",
		c.GetGoModule(),
		"internal/pkg/generated/configs")
}

func (c *Description) GetJsServicesRelPath() string {
	return "js-services"
}

func (c *Description) GetJsServicesAbsPath() string {
	return path.Join(c.BasePath, c.GetJsServicesRelPath())
}

func (c *Description) GetJsServiceRelPath(serviceName string) string {
	return path.Join(c.GetJsServicesRelPath(), serviceName)
}

func (c *Description) GetJsServiceAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsServiceRelPath(serviceName))
}

func (c *Description) GetJsPackageJsonRelPath() string {
	return path.Join(c.GetJsServicesRelPath(), "package.json")
}

func (c *Description) GetJsPackageJsonAbsPath() string {
	return path.Join(c.BasePath, c.GetJsPackageJsonRelPath())
}

func (c *Description) GetJsServicePackageJsonRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "package.json")
}

func (c *Description) GetJsServicePackageJsonAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsServicePackageJsonRelPath(serviceName))
}

func (c *Description) GetJsServiceYarnLockRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "yarn.lock")
}

func (c *Description) GetJsServiceYarnLockAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsServiceYarnLockRelPath(serviceName))
}

func (c *Description) GetJsServiceNuxtConfigRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "nuxt.config.js")
}

func (c *Description) GetJsServiceNuxtConfigAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsServiceNuxtConfigRelPath(serviceName))
}

func (c *Description) GetJsDockerfileRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "Dockerfile")
}

func (c *Description) GetJsDockerfileAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsDockerfileRelPath(serviceName))
}

func (c *Description) GetJsPagesRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "pages")
}

func (c *Description) GetJsPagesAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsPagesRelPath(serviceName))
}

func (c *Description) GetJsIndexRelPath(serviceName string) string {
	return path.Join(c.GetJsPagesRelPath(serviceName), "index.vue")
}

func (c *Description) GetJsIndexAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsIndexRelPath(serviceName))
}

func (c *Description) GetJsComponentsRelPath(serviceName string) string {
	return path.Join(c.GetJsServiceRelPath(serviceName), "components")
}

func (c *Description) GetJsComponentsAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsComponentsRelPath(serviceName))
}

func (c *Description) GetJsSampleVueRelPath(serviceName string) string {
	return path.Join(c.GetJsComponentsRelPath(serviceName), "sample.vue")
}

func (c *Description) GetJsSampleVueAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsSampleVueRelPath(serviceName))
}

func (c *Description) GetJsServicesPath() string {
	return path.Join(c.BasePath, "js-services")
}


func (c *Description) GetGoServicesRelPath() string {
	return "go-services"
}

func (c *Description) GetGoServiceGeneratedCoreRelPath(serviceName string) string {
	return path.Join(c.GetGoServicesRelPath(), "internal", serviceName, "generated/core")
}

func (c *Description) GetGoServicesAbsPath() string {
	return path.Join(c.BasePath, c.GetGoServicesRelPath())
}

func (c *Description) GetGoModRelPath() string {
	return path.Join(c.GetGoServicesRelPath(), "go.mod")
}

func (c *Description) GetGoModAbsPath() string {
	return path.Join(c.BasePath, c.GetGoModRelPath())
}

func (c *Description) GetGoSumRelPath() string {
	return path.Join(c.GetGoServicesRelPath(), "go.sum")
}

func (c *Description) GetGoSumAbsPath() string {
	return path.Join(c.BasePath, c.GetGoSumRelPath())
}


func (c *Description) GetPythonServicesRelPath() string {
	return "py-services"
}

func (c *Description) GetPythonServicesAbsPath() string {
	return path.Join(c.BasePath, c.GetPythonServicesRelPath())
}

func (c *Description) GetPythonServicesLibrariesGeneratedRelPath() string {
	return path.Join(c.GetPythonServicesRelPath(), "libraries/generated")
}

func (c *Description) GetPythonServicesLibrariesGeneratedConfigsRelPath() string {
	return path.Join(c.GetPythonServicesLibrariesGeneratedRelPath(), "configs")
}

func (c *Description) GetPythonServicesLibrariesGeneratedLogsRelPath() string {
	return path.Join(c.GetPythonServicesLibrariesGeneratedRelPath(), "logs")
}

func (c *Description) GetPythonServicesLibrariesGeneratedMetricsRelPath() string {
	return path.Join(c.GetPythonServicesLibrariesGeneratedRelPath(), "metrics")
}

func (c *Description) GetPythonServicesLibrariesGeneratedAbsPath() string {
	return path.Join(c.BasePath, c.GetPythonServicesLibrariesGeneratedRelPath())
}

func (c *Description) GetPythonServicesLibrariesGeneratedConfigsAbsPath() string {
	return path.Join(c.BasePath, c.GetPythonServicesLibrariesGeneratedConfigsRelPath())
}

func (c *Description) GetPythonServicesLibrariesGeneratedLogsAbsPath() string {
	return path.Join(c.BasePath, c.GetPythonServicesLibrariesGeneratedLogsRelPath())
}

func (c *Description) GetPythonServicesLibrariesGeneratedMetricsAbsPath() string {
	return path.Join(c.BasePath, c.GetPythonServicesLibrariesGeneratedMetricsRelPath())
}


func (c *Description) GetPythonGeneratedRelPath(serviceName string) string {
	return path.Join(mifyconfig.PythonServicesRoot, serviceName, "generated")
}

func (c *Description) GetPythonServiceRelPath(serviceName string) string {
	return path.Join(c.GetPythonServicesRelPath(), serviceName)
}

func (c *Description) GetPythonServiceSubAbsPath(serviceName string, filename string) string {
	return path.Join(c.GetPythonServicesAbsPath(), serviceName, filename)
}

func (c *Description) GetPythonAppRelPath(serviceName string) string {
	return path.Join(c.GetPythonServicesRelPath(), serviceName, "app")
}

func (c *Description) GetPythonGeneratedAppPath(serviceName string) string {
	return path.Join(c.GetPythonServicesAbsPath(), serviceName, "generated/app")
}

func (c *Description) GetPythonGeneratedAppRelPath(serviceName string) string {
	return path.Join(c.GetPythonGeneratedRelPath(serviceName), "app")
}

func (c *Description) GetPythonGeneratedAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetPythonGeneratedRelPath(serviceName))
}

func (c *Description) GetPythonAppSubRelPath(serviceName string, fileName string) string {
	return path.Join(c.GetPythonAppRelPath(serviceName), fileName)
}

func (c *Description) GetPythonAppSubAbsPath(serviceName string, fileName string) string {
	return path.Join(c.BasePath, c.GetPythonAppSubRelPath(serviceName, fileName))
}

func (c *Description) GetPythonServiceGeneratedCoreRelPath(serviceName string) string {
	return path.Join(c.GetPythonServicesRelPath(), serviceName, "generated/core")
}

func (c *Description) GetPythonServiceGeneratedOpenAPIRelPath(serviceName string) string {
	return path.Join(c.GetPythonServicesRelPath(), serviceName, "generated/openapi")
}


func (c *Description) GetDevRunnerRelPath() string {
	return c.GetCmdRelPath(DevRunnerName)
}

func (c *Description) GetDevRunnerAbsPath() string {
	return path.Join(c.GetCmdAbsPath(DevRunnerName))
}

func (c *Description) GetDevRunnerMainRelPath() string {
	return path.Join(c.GetDevRunnerRelPath(), "main.go")
}

func (c *Description) GetDevRunnerMainAbsPath() string {
	return path.Join(c.BasePath, c.GetDevRunnerMainRelPath())
}

func (c *Description) GetServicesAbsPath(lang mifyconfig.ServiceLanguage) (string, error) {
	switch lang {
	case mifyconfig.ServiceLanguageGo:
		return c.GetGoServicesAbsPath(), nil
	case mifyconfig.ServiceLanguagePython:
		return c.GetPythonServicesAbsPath(), nil
	case mifyconfig.ServiceLanguageJs:
		return c.GetJsServicesPath(), nil
	}
	return "", ErrUnsupportedLanguage
}

func (c *Description) GetDockerfileAbsPath(serviceName string, lang mifyconfig.ServiceLanguage) (string, error) {
	if !c.HasService(serviceName) {
		return "", ErrNoSuchService
	}
	switch lang {
	case mifyconfig.ServiceLanguageGo:
		return path.Join(c.GetCmdAbsPath(serviceName), "Dockerfile"), nil
	case mifyconfig.ServiceLanguagePython:
		return path.Join(c.GetPythonServicesAbsPath(), serviceName, "Dockerfile"), nil
	case mifyconfig.ServiceLanguageJs:
		return path.Join(c.GetJsServicesPath(), serviceName, "Dockerfile"), nil
	}
	return "", ErrUnsupportedLanguage
}

func (c *Description) GetServiceGeneratedAPIRelPath(serviceName string, language mifyconfig.ServiceLanguage) (string, error) {
	switch language {
	case mifyconfig.ServiceLanguageGo:
		return mifyconfig.GoServicesRoot + "/internal/" + serviceName, nil
	case mifyconfig.ServiceLanguageJs:
		return mifyconfig.JsServicesRoot + "/" + serviceName, nil
	case mifyconfig.ServiceLanguagePython:
		return mifyconfig.PythonServicesRoot + "/" + serviceName, nil
	}
	return "", ErrUnsupportedLanguage
}


func (c *Description) GetCmdAbsPath(serviceName string) string {
	return path.Join(c.GetGoServicesAbsPath(), "cmd", serviceName)
}

func (c *Description) GetCmdRelPath(serviceName string) string {
	return path.Join(c.GetGoServicesRelPath(), "cmd", serviceName)
}

func (c *Description) GetGeneratedRelPath(serviceName string) string {
	return path.Join(mifyconfig.GoServicesRoot, "internal", serviceName, "generated")
}

func (c *Description) GetGeneratedAppRelPath(serviceName string) string {
	return path.Join(c.GetGeneratedRelPath(serviceName), "app")
}

func (c *Description) GetGeneratedAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetGeneratedRelPath(serviceName))
}

func (c *Description) GetJsGeneratedRelPath(serviceName string) string {
	return path.Join(mifyconfig.JsServicesRoot, serviceName, "generated")
}

func (c *Description) GetJsGeneratedAbsPath(serviceName string) string {
	return path.Join(c.BasePath, c.GetJsGeneratedRelPath(serviceName))
}

func (c *Description) GetGeneratedAppPath(serviceName string) string {
	return path.Join(c.GetGoServicesAbsPath(), "internal", serviceName, "generated/app")
}

func (c *Description) GetAppRelPath(serviceName string) string {
	return path.Join(c.GetGoServicesRelPath(), "internal", serviceName, "app")
}

func (c Description) GetGoPostgresConfigRelPath() string {
	return path.Join(c.GetGoServicesRelPath(), "internal/pkg/generated/postgres")
}

func (c Description) GetGoPostgresConfigAbsPath() string {
	return path.Join(c.BasePath, c.GetGoPostgresConfigRelPath())
}

// User app
func (c *Description) GetAppSubRelPath(serviceName string, fileName string) string {
	return path.Join(c.GetAppRelPath(serviceName), fileName)
}

func (c *Description) GetAppSubAbsPath(serviceName string, fileName string) string {
	return path.Join(c.BasePath, c.GetAppSubRelPath(serviceName, fileName))
}

// Name which can be used in generated go code
func (c GoService) GetSafeName() string {
	return util.ToSafeGoVariableName(c.Name)
}

// Mify cache

func (c *Description) GetCacheDirectory() string {
	return filepath.Join(c.BasePath, ".mify")
}

func (c *Description) GetLogsDirectory() string {
	return filepath.Join(c.GetCacheDirectory(), "logs")
}

func (c *Description) GetServiceCacheDirectory(serviceName string) string {
	return filepath.Join(c.GetCacheDirectory(), TmpSubdir, serviceName)
}

func (c *Description) HasApi(serviceName string) bool {
	if _, err := os.Stat(c.GetApiSchemaAbsPath(serviceName, MainApiSchemaName)); os.IsNotExist(err) {
		return false
	}

	return true
}

// Postgres

func (c *Description) GetMigrationsDirectory(databaseName string, lang mifyconfig.ServiceLanguage) (string, error) {
	switch lang {
	case mifyconfig.ServiceLanguageGo:
		return filepath.Join(c.GetGoServicesAbsPath(), "migrations", databaseName), nil
	}
	return "", ErrUnsupportedLanguage
}
