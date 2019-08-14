# checkbot
Checkbot is able to run custom bash script in a container running on OpenShift. These scripts can check funcitonality and compliance settings in your cluster.

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

## OpenShift
```
# create new project
oc new-project checkbot

# setup build
oc new-build https://github.com/tobiasdenzler/checkbot

# start build
oc start-build -F checkbot

# setup
oc apply -f openshift/setup

# create configmaps
oc create configmap scripts-compliance --from-file=scripts/compliance
oc create configmap scripts-operation --from-file=scripts/operation
```

## Features

* Reload endpoint for loading scripts on the fly
* Dashboard that lists all checks
