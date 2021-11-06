{{- .Workspace.TplHeader}}

package core

import (
	"os"

	"go.uber.org/zap"
	"{{.GoModule}}/internal/{{.ServiceName}}/app"
)

type MifyServiceContext struct {
	serviceName string
	hostname    string

	loggerWrapper *MifyLoggerWrapper
	dynamicConfig *MifyDynamicConfig

	serviceContext *app.ServiceContext
}

func NewMifyServiceContext(serviceName string) (*MifyServiceContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &MifyServiceContext{}, err
	}

	context := &MifyServiceContext{
		serviceName: serviceName,
		hostname:    hostname,
	}

	logger, err := NewMifyLoggerWrapper(context)
	if err != nil {
		return &MifyServiceContext{}, nil
	}
	context.loggerWrapper = logger

	dynamicConfig, err := NewMifyDynamicConfig(context)
	if err != nil {
		return nil, err
	}
	context.dynamicConfig = dynamicConfig

	svcCtx, err := app.NewServiceContext()
	if err != nil {
		return &MifyServiceContext{}, err
	}
	context.serviceContext = svcCtx

	return context, nil
}

func (c *MifyServiceContext) ServiceName() string {
	return c.serviceName
}

func (c *MifyServiceContext) Hostname() string {
	return c.hostname
}

func (c *MifyServiceContext) Logger() *zap.Logger {
	return c.loggerWrapper.Logger()
}

func (c *MifyServiceContext) LoggerFor(component string) *zap.Logger {
	return c.loggerWrapper.LoggerFor(component)
}

func (c *MifyServiceContext) ServiceContext() *app.ServiceContext {
	return c.serviceContext
}
