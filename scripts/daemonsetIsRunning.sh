#!/bin/sh

# Check if all pods from Daemonset are running.

set -eux

SCHEDULED=$(oc get ds/checker-daemonset -o json | jq .status.desiredNumberScheduled)
AVAILABLE=$(oc get ds/checker-daemonset -o json | jq .status.numberAvailable)

if [ $SCHEDULED -eq $AVAILABLE ]
then
  # check passed
  echo "ok"
  exit 0
fi

# check failed
echo "we need $SCHEDULED pods but $AVAILABLE are available"
exit 1