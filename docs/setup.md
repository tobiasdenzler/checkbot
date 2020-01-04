# Setup

## Development

Run the server locally:

```
make run
```

There are some flags you can use to configure the application:

Flag | Description | Values |
--- | --- | ---
scriptBase | Base path for the check scripts | e.g. scripts
metricsPrefix | Prefix for all metrics | e.g. checkbot 
logLevel | Log level for application | error &#124; warn &#124; info &#124; debug &#124; trace 
reloadPassword | Password for reload endpoint | e.g. secret 
enableSandbox | Enable debugging sandbox | true &#124; false 

Run the tests:

```
make tests
```

For local development you can create new server certs like this:
```
openssl genrsa -out tls.key 2048
openssl req -new -x509 -sha256 -key tls.key -out tls.crt -days 3650
```

## Docker

You can get a prebuilt image with

```
docker pull tobiasdenzler/checkbot:latest
```

or use one of the [released versions](https://hub.docker.com/repository/docker/tobiasdenzler/checkbot/tags?page=1).

## Openshift

There is some predefined configuration you can use to setup checkbot on Openshift. TLS is implemented using [Service Serving Certificate Secrets](https://docs.openshift.com/container-platform/3.11/dev_guide/secrets.html#service-serving-certificate-secrets).

```
# create new namespace
oc apply -f installation/all/namespace.yaml

# create configmaps
oc -n checkbot create configmap scripts-compliance --from-file=scripts/compliance
oc -n checkbot create configmap scripts-operation --from-file=scripts/operation

# create deployment
oc -n checkbot apply -f installation/all
oc -n checkbot apply -f installation/openshift

```

## Kubernetes

To secure the access to checkbot using TLS you need to generate valid certificates first, please check [Manage TLS Certificates in a Cluster](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/) on how to do that. If you have approved certificates you can add them as a secret to your namespace:

```
kubectl -n checkbot create secret tls checkbot-certs --cert=tls.crt --key=tls.key
```

Now you can use the following scripts to setup checkbot on Kubernetes:

```
# create new namespace
kubectl apply -f installation/all/namespace.yaml

# create configmaps
kubectl -n checkbot create configmap scripts-compliance --from-file=scripts/compliance
kubectl -n checkbot create configmap scripts-operation --from-file=scripts/operation

# create deployment
kubectl -n checkbot apply -f installation/all
kubectl -n checkbot apply -f installation/kubernetes
```
