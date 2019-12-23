# checkbot

![](https://github.com/tobiasdenzler/checkbot/workflows/build/badge.svg)

Checkbot is able to run custom bash script in a container running on OpenShift. These scripts can check functionality and compliance settings in your cluster and will expose the result as Prometheus metrics.

![Checkbot Overview](checkbot_overview.png)

## How it Works

Checks are written in Bash and need to be saved as .sh files. Each check will provide results in form of one [type of metric](https://prometheus.io/docs/concepts/metric_types/).

A check must contain some metadata for registering the check. Metadata is written as comment and need to contain the following information:

* ACTIVE: Is the check currently active (true|false)
* TYPE: The type of the metric (Gauge)
* HELP: Description of the metric
* INTERVAL: Number of seconds between runs of the check

The return values need to follow a predefined format:
```
value|label1=value1,label2=value2
```
It is also possible to return multiple lines. But be sure that you provide the same labels on each line otherwise it would not be a valid metric.

Example:

```
#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 60

set -eux

# file1 contains all projects
oc get project --no-headers | awk '{print $1}' | sort > /tmp/file1

# file2 contains all quotas
oc get quota --all-namespaces --no-headers | awk '{print $1}' | sort| uniq > /tmp/file2

# result contains projects without quotas
comm -3 /tmp/file1 /tmp/file2 > /tmp/result

# looping through results
while IFS="" read -r p || [ -n "$p" ]
do
  printf '1|project=%s\n' "$p"
done < /tmp/result

exit 0
```

This script will produce results like the following:

```
1|project=grafana
1|project=kube-dns
1|project=test
```
Checkbot will then read the result and convert it into the appropriate Prometheus metric:

```
# HELP checkbot_missing_quota_on_project_total Check if all projects have quotas defined.
# TYPE checkbot_missing_quota_on_project_total gauge
checkbot_missing_quota_on_project_total{project="grafana"} 1
checkbot_missing_quota_on_project_total{project="kube-dns"} 1
checkbot_missing_quota_on_project_total{project="test"} 1
```

## Local Development

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

### Docker

Use Docker to build the image locally

```
# build images
docker build -t checkbot .

# run image
docker run checkbot
docker run -it checkbot /bin/bash
```


### Minishift

Use Minishift for integration tests with OpenShift:

```
# install addons, check https://github.com/minishift/minishift-addons
minishift addon enable admin-user
minishift addon apply admin-user
minishift addon enable registry-route
minishift addon apply registry-route
minishift addon install prometheus
minishift addon enable prometheus
minishift addon apply prometheus --addon-env namespace=kube-system
minishift addon install management-infra
minishift addon enable management-infra
minishift addon apply management-infra
minishift addon install grafana
minishift addon enable grafana
minishift addon apply grafana --addon-env namespace=grafana

# starting minishift
minishift start --v 5 --cpus=4

# login
oc login -u system:admin
```


## Cluster Setup

To operate checkbot on your Openshift or Kubernetes cluster the following steps might be helpful.

### OpenShift
```
# create new project
oc new-project checkbot

# setup build
oc new-build https://github.com/tobiasdenzler/checkbot

# start build
oc start-build -F checkbot

# create configmaps
oc create configmap scripts-compliance --from-file=scripts/compliance
oc create configmap scripts-operation --from-file=scripts/operation

# setup
oc apply -f openshift/setup

```

### Prometheus

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

## Hints

### Reload

If you change the scripts in your configmap you can use the reload endpoint to reload all scripts:
```
curl -k -X POST -u admin:admin https://localhost:4444/reload
```
Default values for authentication using basic auth are admin/admin. The default password for the reload endpoint can be changed using the --reloadPassword flag.

### Generate server certificates

For local development you can create new server certs like this:
```
openssl genrsa -out tls.key 2048
openssl req -new -x509 -sha256 -key tls.key -out tls.crt -days 3650
```
For OpenShift you can use the service certificates.

### Sandbox

There is a sandbox you can use to test and debug your check scripts. You have to enable this feature by using the -enableSandbox=true flag.

> Be aware that the sandbox is able to execute any script you paste and therefore is able to control its container or your local environment.

### Lastrun

To check if your scripts have run successfully you can use the (internal) metric lastrun_info. This metric will provide information about the last run of each check:

```
checkbot_lastrun_info{interval="60",name="checkbot_missing_quota_on_project_total",offset="22",success="true",type="Gauge"} 1.57699768e+09
checkbot_lastrun_info{interval="60",name="checkbot_modified_scc_reconcile",offset="12",success="true",type="Gauge"} 1.576997641e+09
```

Note:  Offset is the number of second that is used to randomly delay the execution of the script. To get the time of the next run you can add the interval and the offset to the current time.
