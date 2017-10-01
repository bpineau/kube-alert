package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/controllers"
	"github.com/bpineau/kube-alert/pkg/handlers"
	"github.com/bpineau/kube-alert/pkg/health"
)

func Run(config *config.AlertConfig) {
	for _, controller := range controllers.Controllers {
		go controller.Start(config, handlers.Handlers[controller.HandlerName()])
	}

	go health.HealthCheckServe(config)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
