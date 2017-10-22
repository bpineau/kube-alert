package handlers

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/handlers/cs"
	"github.com/bpineau/kube-alert/pkg/handlers/node"
	"github.com/bpineau/kube-alert/pkg/handlers/pod"
)

// Handler reacts and analyze controllers provided events.
type Handler interface {
	Init(c *config.AlertConfig) error

	// ObjectCreated should return true and an empty string
	// if the objet is healthy, false and a message otherwise.
	ObjectCreated(obj interface{}) (bool, string)

	// ObjectDeleted should return true and an empty string
	// if the objet is healthy, false and a message otherwise.
	ObjectDeleted(obj interface{}) (bool, string)
}

// Handlers map all known handlers
var Handlers = map[string]Handler{
	"cs":   &cs.Handler{},
	"pod":  &pod.Handler{},
	"node": &node.Handler{},
}
