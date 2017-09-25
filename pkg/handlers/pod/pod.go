package pod

import (
	"fmt"
	"github.com/bpineau/kube-alert/config"
	"k8s.io/client-go/pkg/api/v1"
	"time"
)

var (
	minAgeMinutes time.Duration = 15
)

type PodHandler struct {
	conf *config.AlertConfig
}

func (h *PodHandler) Init(c *config.AlertConfig) error {
	c.Logger.Info("pod handler initialized")
	h.conf = c
	return nil
}

func (h *PodHandler) ObjectCreated(obj interface{}) (bool, string) {
	pod, _ := obj.(*v1.Pod)

	// ignore recent pods
	if pod.Status.StartTime == nil {
		return true, ""
	}
	if time.Now().Add(-minAgeMinutes * time.Minute).Before(pod.Status.StartTime.Time) {
		return true, ""
	}

	healthy, reason := h.checkPodHealthy(pod)
	if !healthy {
		return false, fmt.Sprintf("%s/%s is unhealthy: %s", pod.Namespace, pod.Name, reason)
	}

	return true, ""
}

func (h *PodHandler) ObjectDeleted(obj interface{}) (bool, string) {
	return true, ""
}

func (h *PodHandler) checkPodHealthy(pod *v1.Pod) (bool, string) {

	if pod.Status.Phase == "Failed" {
		return false, "pod in Failed state: " + extractContainersErrors(pod)
	}

	if pod.Status.Phase == "Pending" {
		return false, "pod remains on Pending state: " + extractContainersErrors(pod)
	}

	for _, container := range pod.Status.ContainerStatuses {
		if int(container.RestartCount) == 0 || container.State.Running == nil {
			continue
		}
		if time.Now().Add(-minAgeMinutes * time.Minute).Before(container.State.Running.StartedAt.Time) {
			return false, fmt.Sprintf("container restarted %d times, last at %s",
				container.RestartCount, container.State.Running.StartedAt)
		}
	}

	return true, ""
}

func extractContainersErrors(pod *v1.Pod) string {
	for _, container := range pod.Status.ContainerStatuses {
		if container.Ready {
			continue
		}
		if container.State.Waiting != nil {
			return container.State.Waiting.Reason
		}
		if container.State.Terminated != nil {
			return container.State.Terminated.Reason
		}
	}

	return ""
}
