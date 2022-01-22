package app

import (
	"net/http"

	"github.com/chebykinn/mify-deploy/go_services/internal/k8s-control/generated/core"
)

type routerConfig struct {
	Middlewares []func(http.Handler) http.Handler
}

func NewRouterConfig() *routerConfig {
	return &routerConfig {
		Middlewares: []func(http.Handler) http.Handler {
		// Add your middlewares here
		},
	}
}

type ServiceContext struct {
	// Append your dependencies here
}

func NewServiceContext(ctx *core.MifyServiceContext) (*ServiceContext, error) {
	// Here you can do your custom service initialization, prepare dependencies
	context := &ServiceContext{
		// Here you can initialize your dependencies
	}
	return context, nil
}
