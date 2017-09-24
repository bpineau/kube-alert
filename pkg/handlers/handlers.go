package handlers

import (
	"github.com/bpineau/kube-alert/config"
)

type Handler interface {
	Init(c *config.AlertConfig) error
	ObjectCreated(obj interface{}) (bool, string)
	ObjectDeleted(obj interface{}) (bool, string)
}
