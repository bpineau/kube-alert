package log

import (
	"github.com/bpineau/kube-alert/config"
)

type LogNotifier struct {
}

func (l *LogNotifier) Notify(c *config.AlertConfig, title string, msg string) error {
	c.Logger.Infof("%s: %s", title, msg)
	return nil
}
