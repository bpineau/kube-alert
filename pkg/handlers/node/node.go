package node

import (
	"fmt"

	"github.com/bpineau/kube-alert/config"
	"k8s.io/client-go/pkg/api/v1"
)

type NodeHandler struct {
	conf *config.AlertConfig
}

var knownBadConditions = map[string]bool{
	"OutOfDisk":          true,
	"MemoryPressure":     true,
	"DiskPressure":       true,
	"NetworkUnavailable": true,
}

func (n *NodeHandler) Init(c *config.AlertConfig) error {
	c.Logger.Info("node handler initialized")
	n.conf = c
	return nil
}

func (n *NodeHandler) ObjectCreated(obj interface{}) (bool, string) {
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

func (n *NodeHandler) ObjectDeleted(obj interface{}) (bool, string) {
	return true, ""
}
