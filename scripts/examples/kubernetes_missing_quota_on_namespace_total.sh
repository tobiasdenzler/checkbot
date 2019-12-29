#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check if all namespaces have resourcequotas defined.
# INTERVAL 60

set -eu

# file1 contains all projects
kubectl get namespace --no-headers | awk '{print $1}' | sort > /tmp/file1

# file2 contains all quotas
kubectl get resourcequota --all-namespaces --no-headers | awk '{print $1}' | sort| uniq > /tmp/file2

# result contains projects without quotas
comm -3 /tmp/file1 /tmp/file2 > /tmp/result

# looping through results and print them out
while IFS="" read -r p || [ -n "$p" ]
do
  printf '1|project=%s\n' "$p"
done < /tmp/result

exit 0