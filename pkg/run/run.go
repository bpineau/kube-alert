package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/controllers"
	"github.com/bpineau/kube-alert/pkg/controllers/cs"
	"github.com/bpineau/kube-alert/pkg/controllers/node"
	"github.com/bpineau/kube-alert/pkg/controllers/pod"
	"github.com/bpineau/kube-alert/pkg/handlers"
	"github.com/bpineau/kube-alert/pkg/health"
)

var Controllers = []controllers.Controller{
	&cs.CsController{},
	&pod.PodController{},
	&node.NodeController{},
}

func Run(config *config.AlertConfig) {
	for _, c := range Controllers {
		go c.Init(config, handlers.Handlers[c.HandlerName()]).Start()
	}

	go health.HealthCheckServe(config)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
