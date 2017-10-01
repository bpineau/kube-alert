package handlers

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/handlers/pod"
	"github.com/bpineau/kube-alert/pkg/handlers/cs"
)

type Handler interface {
	Init(c *config.AlertConfig) error
	ObjectCreated(obj interface{}) (bool, string)
	ObjectDeleted(obj interface{}) (bool, string)
}

var Handlers = map[string]Handler{
	"pod": &pod.PodHandler{},
	"cs": &cs.CsHandler{},
}
