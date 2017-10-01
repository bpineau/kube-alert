package pod

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

type PodController struct {
	// https://golang.org/doc/effective_go.html#embedding
	controllers.CommonController
}

func (c *PodController) HandlerName() string {
	return "pod"
}

func (c *PodController) Init(conf *config.AlertConfig, handler handlers.Handler) controllers.Controller {
	c.CommonController = controllers.CommonController{}
	c.Conf = conf
	c.Handler = handler

	client := c.Conf.ClientSet
	c.Name = "pod"
	c.ObjType = &v1.Pod{}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
		},
	}

	return c
}
