#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all pods from Daemonset are running.
# INTERVAL 10

set -eux

SCHEDULED=$(oc get ds/checkbot-daemonset -o json | jq .status.desiredNumberScheduled)
AVAILABLE=$(oc get ds/checkbot-daemonset -o json | jq .status.numberAvailable)

DIFF="$(($SCHEDULED-$AVAILABLE))"

# The return value of the script should contain a value and a map of labels
# value(int)|label1:value1,label2:value2
echo "$DIFF"
exit 0