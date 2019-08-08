# openshift-healthchecker
Provides different HTTP healthchecks of OpenShift components.

## Start server

```
go run ./cmd/server
```

## Checks

- check/daemonsetIsRunning
- check/projectHasQuota

## Minishift

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

## API
```
curl -ik -X GET -H "Authorization: Bearer sprO8RxGNh8swwy_pmYmfT_GaHKL3EeKfU2ASovXDec" https://192.168.42.28:8443/oapi/v1

curl -ik -X GET -H 'Accept: application/json' -H "Authorization: Bearer sprO8RxGNh8swwy_pmYmfT_GaHKL3EeKfU2ASovXDec" https://192.168.42.28:8443/apis/proect.openshift.io/v1/projects
```

## Docker
```
# build images
docker build -t openshift-healthchecker .

# run image
docker run openshift-healthchecker
```

## OpenShift
```
# build image
oc start-build -F openshift-healthchecker
```
