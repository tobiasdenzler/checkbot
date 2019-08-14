#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 60

set -eux

PROJECTS=$(oc get project --no-headers | wc -l)
QUOTAS=$(oc get quota --all-namespaces --no-headers | wc -l)

if [ $PROJECTS -eq $QUOTAS ]
then
  # check passed
  echo "ok"
  exit 0
fi

# check failed
echo "found $PROJECTS project and $QUOTAS quotas"
exit 1