package log

import (
	"github.com/bpineau/kube-alert/config"
)

// Notifier implements notifiers.Notifier
type Notifier struct {
}

// Notify sends notification to the configured logrus logger
func (l *Notifier) Notify(c *config.AlertConfig, title string, msg string) error {
	c.Logger.Infof("%s: %s", title, msg)
	return nil
}
