package notifiers

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/notifiers/datadog"
	"github.com/bpineau/kube-alert/pkg/notifiers/log"
)

type Notifier interface {
	Notify(c *config.AlertConfig, title string, msg string) error
}

var Notifiers = []Notifier{
	&log.LogNotifier{},
	&datadog.DdNotifier{},
}
