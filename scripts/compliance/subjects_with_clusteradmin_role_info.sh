#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check all subjects with cluster-admin role.
# INTERVAL 3600

set -eux

# Retrieve all subjects with cluster-admin role
oc get clusterrolebinding -o json \
  | jq '.items[] | select(.roleRef.name | startswith("cluster-admin")) | .subjects[] | "1|subject="+.kind+","+"name="+.name' \
  | tr -d "\""

exit 0
