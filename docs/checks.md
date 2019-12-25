# Writing Checks

Checks are written in Bash and need to be saved as .sh files. Each check will provide results in form of one [type of metric](https://prometheus.io/docs/concepts/metric_types/).

## Configuration

A check must contain some metadata for registering the check. Metadata is written as comment and need to contain the following information:

* ACTIVE: Is the check currently active (true|false)
* TYPE: The type of the metric (Gauge)
* HELP: Description of the metric
* INTERVAL: Number of seconds between runs of the check

## Return Values

The return values need to follow a predefined format:
```
value|label1=value1,label2=value2
```
It is also possible to return multiple lines. But be sure that you provide the same labels on each line otherwise it would not be a valid metric.

## Example

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

## Reload

If you change the scripts in your configmap you can use the reload endpoint to reload all scripts:
```
curl -k -X POST -u admin:admin https://localhost:4444/reload
```
Default values for authentication using basic auth are admin/admin. The default password for the reload endpoint can be changed using the --reloadPassword flag.

## Sandbox

There is a sandbox you can use to test and debug your check scripts. You have to enable this feature by using the -enableSandbox=true flag.

> Be aware that the sandbox is able to execute any script you paste and therefore is able to control its container or your local environment.

