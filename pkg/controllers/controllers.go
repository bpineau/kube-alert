package controllers

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/controllers/pod"
	"github.com/bpineau/kube-alert/pkg/handlers"
)

type Controller interface {
	Start(c *config.AlertConfig, handler handlers.Handler)
	HandlerName() string
}

var Controllers = []Controller{
	&pod.PodController{},
}
