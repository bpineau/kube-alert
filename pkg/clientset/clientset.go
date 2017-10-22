package clientset

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Ensure we have GCP auth method linked in
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func buildConfig(apiserver string, kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".kube/config")); err == nil {
			kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube/config")
		}
	}

	if apiserver != "" || kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	}

	return rest.InClusterConfig()
}

// NewClientSet create a clientset for the optional apiserver or kubeconfig configs,
// defaulting to the automatic, in cluster settings.
func NewClientSet(apiserver string, kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := buildConfig(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
