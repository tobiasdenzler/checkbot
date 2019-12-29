#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check modified scc subjects
# INTERVAL 3600

set -eux

# Retrieve all subjects with cluster-admin role
oc adm policy reconcile-sccs --confirm=false -ojson  \
  | jq '.items | .[] | "1|name="+.metadata.name '\
  | tr -d "\""

exit 0
