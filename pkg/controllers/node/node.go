package node

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/controllers"
	"github.com/bpineau/kube-alert/pkg/handlers"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

type NodeController struct {
	// https://golang.org/doc/effective_go.html#embedding
	controllers.CommonController
}

func (c *NodeController) HandlerName() string {
	return "node"
}

func (c *NodeController) Init(conf *config.AlertConfig, handler handlers.Handler) controllers.Controller {
	c.CommonController = controllers.CommonController{
		Conf:    conf,
		Handler: handler,
		Name:    "node",
	}

	client := c.Conf.ClientSet
	c.ObjType = &v1.Node{}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Nodes().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Nodes().Watch(options)
		},
	}

	return c
}
