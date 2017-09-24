package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/controllers/pod"
	pod_handler "github.com/bpineau/kube-alert/pkg/handlers/pod"
	"github.com/bpineau/kube-alert/pkg/notifiers"
)

func Run(config *config.AlertConfig) {
	go pod.Start(config, new(pod_handler.PodHandler), notifiers.Notifiers)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}