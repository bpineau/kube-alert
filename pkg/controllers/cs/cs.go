package cs

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

type CsController struct {
	// https://golang.org/doc/effective_go.html#embedding
	controllers.CommonController
}

func (c *CsController) HandlerName() string {
	return "cs"
}

func (c *CsController) Init(conf *config.AlertConfig, handler handlers.Handler) controllers.Controller {
	c.CommonController = controllers.CommonController{}
	c.Conf = conf
	c.Handler = handler

	client := c.Conf.ClientSet
	c.Name = "componentstatus"
	c.ObjType = &v1.ComponentStatus{}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().ComponentStatuses().List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			//Apparently this resource does not support watch :
			//   "Failed to watch *v1.ComponentStatus: the server does not allow
			//    this method on the requested resource (get componentstatuses)"
			//So we'll only see errors at resync intervals.
			//return client.CoreV1().ComponentStatuses().Watch(options)
			fw := watch.NewFake()
			return fw, nil
		},
	}

	return c
}
