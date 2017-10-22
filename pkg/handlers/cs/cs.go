package cs

import (
	"fmt"

	"github.com/bpineau/kube-alert/config"
	"k8s.io/client-go/pkg/api/v1"
)

// Handler implements handlers.Handler
type Handler struct {
	conf *config.AlertConfig
}

// Init initialize a new cs handler
func (h *Handler) Init(c *config.AlertConfig) error {
	c.Logger.Info("componentstatus handler initialized")
	h.conf = c
	return nil
}

// ObjectCreated inspect a cs health
func (h *Handler) ObjectCreated(obj interface{}) (bool, string) {
	cs, _ := obj.(*v1.ComponentStatus)

	if cs == nil || cs.Conditions == nil {
		return true, ""
	}

	for _, c := range cs.Conditions {
		if c.Type != "Healthy" {
			continue
		}
		if c.Status != "True" {
			return false, fmt.Sprintf("%s is unhealthy: %s", cs.Name, c.Message)
		}
	}

	return true, ""
}

// ObjectDeleted is notified on cs deletion (won't happen often;)
func (h *Handler) ObjectDeleted(obj interface{}) (bool, string) {
	return true, ""
}
