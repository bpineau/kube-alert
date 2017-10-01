package datadog

import (
	"github.com/bpineau/kube-alert/config"
	"github.com/zorkian/go-datadog-api"
)

type DdNotifier struct {
}

var tags = []string{
	"context:kubernetes",
	"origin:kube-alert",
}

func (l *DdNotifier) Notify(c *config.AlertConfig, title string, msg string) error {
	if c.DdApiKey == "" {
		c.Logger.Debug("Omitting datadog notification, api key missing")
		return nil
	}

	client := datadog.NewClient(c.DdApiKey, c.DdAppKey)

	_, err := client.PostEvent(&datadog.Event{
		Title:     &title,
		Text:      &msg,
		AlertType: datadog.String("error"),
		Tags:      tags,
	})

	if err != nil {
		c.Logger.Warning("failed to post to datadog: %s", err)
	}

	return err
}
