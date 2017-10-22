package notifiers

import (
	"fmt"

	"github.com/bpineau/kube-alert/config"
	"github.com/bpineau/kube-alert/pkg/notifiers/datadog"
	"github.com/bpineau/kube-alert/pkg/notifiers/log"
)

// Notifier sends alerts (title + message) to a backend
type Notifier interface {
	Notify(c *config.AlertConfig, title string, msg string) error
}

// Notifiers maps all known notifiers
var Notifiers = []Notifier{
	&log.Notifier{},
	&datadog.Notifier{},
}

// Notify logs events and calls each notifiers' own Notify()
func Notify(c *config.AlertConfig, title string, msg string) {
	ptitle := title
	if c.MsgPrefix != "" {
		ptitle = fmt.Sprintf("%s %s", c.MsgPrefix, title)
	}

	if c.DryRun {
		fmt.Printf("%s: %s\n", ptitle, msg)
		return
	}

	for _, notifier := range Notifiers {
		err := notifier.Notify(c, ptitle, msg)
		if err != nil {
			c.Logger.Warningf("Failed to notify: %s", err)
		}
	}
}
