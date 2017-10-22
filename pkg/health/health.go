package health

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bpineau/kube-alert/config"
)

func healthCheckReply(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		fmt.Printf("Failed to reply to http healtcheck: %s\n", err)
	}
}

// HeartBeatService exposes an http healthcheck handler
func HeartBeatService(c *config.AlertConfig) {
	if c.HealthPort == 0 {
		return
	}
	http.HandleFunc("/health", healthCheckReply)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.HealthPort), nil); err != nil {
		panic(fmt.Sprintf("Failed to start http healtcheck: %s", err))
	}
}
