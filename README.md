# kube-alert

kube-alert watch for pod failures or anomalies, and send alerts accordingly.

Currently support alerting to Datadog and logs (ie. syslog).

## Build

Assuming you have go 1.9 and glide in the path, and GOPATH configured:

```shell
make dep
make build
```

## Usage

The daemon may run as a pod in the cluster, or outside of the cluster.
He should find the Kubernetes api-server automatically, but you can
provide this server's address with "-s" flag, or a full config with "-k".

You can pass configuration values either by command line arguments, or
environment variables, or with a yaml configuration file.

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
  -p, --healthcheck-port int     port for answering healtchecks
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

