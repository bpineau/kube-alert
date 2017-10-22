package node

import (
	"fmt"

	"github.com/bpineau/kube-alert/config"
	"k8s.io/client-go/pkg/api/v1"
)

// Handler implements handlers.Handler
type Handler struct {
	conf *config.AlertConfig
}

var knownBadConditions = map[string]bool{
	"OutOfDisk":          true,
	"MemoryPressure":     true,
	"DiskPressure":       true,
	"NetworkUnavailable": true,
}

// Init initialize a new node handler
func (n *Handler) Init(c *config.AlertConfig) error {
	c.Logger.Info("node handler initialized")
	n.conf = c
	return nil
}

// ObjectCreated inspect a node health
func (n *Handler) ObjectCreated(obj interface{}) (bool, string) {
	node, _ := obj.(*v1.Node)

	for _, c := range node.Status.Conditions {
		if c.Status == "False" {
			continue
		}

		if knownBadConditions[string(c.Type)] {
			return false, fmt.Sprintf("Node %s is unhealthy: %s", node.Name, c.Message)
		}
	}

	return true, ""
}

// ObjectDeleted is notified on node deletion
func (n *Handler) ObjectDeleted(obj interface{}) (bool, string) {
	return true, ""
}
