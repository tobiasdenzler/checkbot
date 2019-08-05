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
