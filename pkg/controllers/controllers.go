package controllers

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/handlers"
	"github.com/bpineau/kube-alert/pkg/notifiers"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var (
	resyncInterval  time.Duration = 300 * time.Second
	maxProcessRetry int           = 6
)

type Controller interface {
	Start()
	Init(c *config.AlertConfig, handler handlers.Handler) Controller
	HandlerName() string
}

type CommonController struct {
	Conf      *config.AlertConfig
	Queue     workqueue.RateLimitingInterface
	Informer  cache.SharedIndexInformer
	Handler   handlers.Handler
	Name      string
	ListWatch cache.ListerWatcher
	ObjType   runtime.Object
}

func (c *CommonController) Start() {
	if err := c.Handler.Init(c.Conf); err != nil {
		c.Conf.Logger.Fatalf("Failed to init %s handler: %s", c.Name, err)
	}

	c.startInformer()
	stopCh := make(chan struct{})
	defer close(stopCh)

	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func (c *CommonController) startInformer() {
	c.Queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	c.Informer = cache.NewSharedIndexInformer(
		c.ListWatch,
		c.ObjType,
		resyncInterval,
		cache.Indexers{},
	)

	c.Informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.Queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				c.Queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.Queue.Add(key)
			}
		},
	})
}

func (c *CommonController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.Queue.ShutDown()

	c.Conf.Logger.Infof("Starting %s controller", c.Name)

	go c.Informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.Informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.Conf.Logger.Infof("%s controller synced and ready", c.Name)

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *CommonController) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *CommonController) processNextItem() bool {
	key, quit := c.Queue.Get()
	if quit {
		return false
	}
	defer c.Queue.Done(key)

	err := c.processItem(key.(string))

	if err == nil {
		// No error, reset the ratelimit counters
		c.Queue.Forget(key)
	} else if c.Queue.NumRequeues(key) < maxProcessRetry {
		c.Conf.Logger.Errorf("Error processing %s (will retry): %v", key, err)
		c.Queue.AddRateLimited(key)
	} else {
		// err != nil and too many retries
		c.Conf.Logger.Errorf("Error processing %s (giving up): %v", key, err)
		c.Queue.Forget(key)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *CommonController) processItem(key string) error {
	obj, exists, err := c.Informer.GetIndexer().GetByKey(key)

	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}

	if !exists {
		c.Handler.ObjectDeleted(obj)
		return nil
	}

	status, msg := c.Handler.ObjectCreated(obj)
	if !status {
		notifiers.Notify(c.Conf, fmt.Sprintf("%s failure", strings.Title(c.Name)), msg)
	}

	return nil
}
