#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 60

set -eu

# file1 contains all projects
oc get project --no-headers | awk '{print $1}' | sort > /tmp/file1

# file2 contains all quotas
oc get quota --all-namespaces --no-headers | awk '{print $1}' | sort| uniq > /tmp/file2

# result contains projects without quotas
comm -3 /tmp/file1 /tmp/file2 > /tmp/result

# looping through results and print them out
while IFS="" read -r p || [ -n "$p" ]
do
  printf '1|project=%s\n' "$p"
done < /tmp/result

exit 0