#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check modified scc subjects
# INTERVAL 60

set -eux

# Retrieve all subjects with cluster-admin role
oc adm policy reconcile-sccs --confirm=false -ojson  \
  | jq '.items | .[] | "1|name="+.metadata.name '\
  | tr -d "\""

exit 0
