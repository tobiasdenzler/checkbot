# Setup

## Development

You can run the checkbot binary locally using go or as container.

### Go

Run the server locally:

```
go run ./cmd/server
```

Run the tests:

```
go test github.com/tobiasdenzler/checkbot/cmd/server -v
```

## Openshift

To operate checkbot on your Openshift cluster the following steps might be helpful. You can get a prebuilt image with

```
docker pull tobiasdenzler/checkbot:latest
```

or use one of the [released versions](https://hub.docker.com/repository/docker/tobiasdenzler/checkbot/tags?page=1).

There is some predefined configuration you can use to setup checkbot on Openshift:

```
# create new project
oc new-project checkbot

# create configmaps
oc create configmap scripts-compliance --from-file=scripts/compliance
oc create configmap scripts-operation --from-file=scripts/operation

# setup
oc apply -f openshift/setup
```

## Prometheus

Use the following snippet to scrape the checkbot metrics:
```
- job_name: checkbot
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    insecure_skip_verify: true
  static_configs:
    - targets: ['checkbot.checkbot.svc.cluster.local:4444']
```

## Generate server certificates

For local development you can create new server certs like this:
```
openssl genrsa -out tls.key 2048
openssl req -new -x509 -sha256 -key tls.key -out tls.crt -days 3650
```
For OpenShift you can use the service certificates.

### Lastrun

To check if your scripts have run successfully you can use the (internal) metric lastrun_info. This metric will provide information about the last run of each check:

```
checkbot_lastrun_info{interval="60",name="checkbot_missing_quota_on_project_total",offset="22",success="true",type="Gauge"} 1.57699768e+09
checkbot_lastrun_info{interval="60",name="checkbot_modified_scc_reconcile",offset="12",success="true",type="Gauge"} 1.576997641e+09
```

Note:  Offset is the number of second that is used to randomly delay the execution of the script. To get the time of the next run you can add the interval and the offset to the current time.
