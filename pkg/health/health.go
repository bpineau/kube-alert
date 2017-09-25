package health

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bpineau/kube-alert/config"
)

func healthCheckReply(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func HealthCheckServe(c *config.AlertConfig) {
	if c.HealthPort == 0 {
		return
	}
	http.HandleFunc("/health", healthCheckReply)
	http.ListenAndServe(fmt.Sprintf(":%d", c.HealthPort), nil)
}
