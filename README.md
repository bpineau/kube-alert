# kube-alert

kube-alert watch for failures in Kubernetes clusters, and send alerts accordingly.

kube-alert monitors:
* Pod failures (unschedulables, error pulling images, crashloop backoff, etc.)
* Pods restarts
* Cluster's components status (issues with etcd, scheduler, or controller-manager daemons)

Support alerting to Datadog and to logs (ie. syslog).

## Build

Assuming you have go 1.9 and glide in the path, and GOPATH configured:

```shell
make deps
make build
```

## Usage

The daemon may run either as a pod, or outside of the Kubernetes cluster.
He should find the Kubernetes api-server automatically (but you can
provide this server's address with "-s" flag, or a kube config with "-k").

You can pass configuration values either by command line arguments, or
environment variables, a yaml configuration file, or a combination or those.

The command line flags are:
```
Usage:
  kube-alert [flags]

Flags:
  -s, --api-server string        kube api server url
  -c, --config string            configuration file (default "/etc/kube-alert/kube-alert.yaml")
  -i, --datadog-api-key string   datadog api key
  -a, --datadog-app-key string   datadog app key
  -d, --dry-run                  dry-run mode
  -p, --healthcheck-port int     port for answering healthchecks
  -h, --help                     help for kube-alert
  -k, --kube-config string       kube config path
  -v, --log-level string         log level (default "debug")
  -l, --log-output string        log output (default "stderr")
  -r, --log-server string        log server (if using syslog)
```

Using an (optional) configuration file:
```yaml
dry-run: false
healthcheck-port: 8080
api-server: http://example.com:8080

log:
  output: "stdout"
  level: "debug"

datadog:
  api-key: xxx
  app-key: xxx
```

The environment variable consumed by kube-alert are option names prefixed
by ```KUBE_ALERT_``` and using underscore instead of dash. Except KUBECONFIG,
used without a prefix (to match kubernetes conventions).
```
env KUBECONFIG=/etc/kube/config \
    KUBE_ALERT_HEALTHCHECK_PORT=8081 \
    KUBE_ALERT_DATADOG.APP_KEY="xxx" \
    KUBE_ALERT_DATADOG.API_KEY="xxx" \
    KUBE_ALERT_DRY_RUN=true \
    KUBE_ALERT_LOG.LEVEL=info \
    KUBE_ALERT_API_SERVER="http://example.com:8080" \
    kube-alert
```
