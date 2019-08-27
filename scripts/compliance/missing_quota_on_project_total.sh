#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 60

set -eux

# file1 contains all projects
oc get project --no-headers | awk '{print $1}' | sort > file1

# file2 contains all quotas
oc get quota --all-namespaces --no-headers | awk '{print $1}' | sort > file2

# result contains projects without quotas
comm -3 file1 file2 > result

# looping through results
while IFS="" read -r p || [ -n "$p" ]
do
  printf '1|project=%s\n' "$p"
done < result

exit 0