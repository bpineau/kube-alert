package cs

import (
	"fmt"

	"github.com/bpineau/kube-alert/config"
	"k8s.io/client-go/pkg/api/v1"
)

type CsHandler struct {
	conf *config.AlertConfig
}

func (h *CsHandler) Init(c *config.AlertConfig) error {
	c.Logger.Info("componentstatus handler initialized")
	h.conf = c
	return nil
}

func (h *CsHandler) ObjectCreated(obj interface{}) (bool, string) {
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

func (h *CsHandler) ObjectDeleted(obj interface{}) (bool, string) {
	return true, ""
}
