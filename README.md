# checkbot
Checkbot is able to run custom bash script in a container running on OpenShift. These scripts can check funcitonality and compliance settings in your cluster and will expose the result as Prometheus metrics.


## Checks

### Scripts

Checks are written in Bash and need to be saved as .sh files.

A check must contain some metadata for registering the check. Metadata is written as comment and need to contain the following information:

* ACTIVE: Is the check currently active (true|false)
* TYPE: The type of the metric (Gauge)
* HELP: Description of the metric
* INTERVAL: Number of seconds between runs of the check

The return values need to follow a predefined format:
```
value|label1=value1,label2=value2
```
It is also possible to return multiple lines.

Example:

```
#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check all subjects with cluster-admin role.
# INTERVAL 60

set -eux

# Retrieve all subjects with cluster-admin role
SUBJECTS=$(oc get clusterrolebinding -o json | jq '.items[] | select(.metadata.name |  startswith("cluster-admin")) | .subjects[] | "subject="+.kind+","+"name="+.name')

for subject in $SUBJECTS
do
    # Return the subjects with cluster-admin role, tr will strip quotes
    echo "1|$subject" | tr -d "\""
done

exit 0
```


### Reload

If you change the scripts in your configmap you can use the reload endpoint to reload all scripts:
```
curl -X POST http://localhost:4444/reload
```


## Development

### Run

Run the server locally:

```
go run ./cmd/server
```

Check -h for runtime configuration.


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


## Setup

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
  static_configs:
    - targets: ['checkbot-checkbot.192.168.42.28.nip.io:443']
```


## Features

* Reload endpoint using POST
* Reload endpoint using authentication
* Run scripts in browser to debug and test
* Support other metric types
* Add more tools (telnet, netcat, etc.)
* Configurable OpenShift CLI version
* Add AWS CLI

