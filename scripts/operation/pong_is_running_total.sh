#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all pods from Daemonset are running.
# INTERVAL 30

set -eux

AVAILABLE=$(oc get ds/pong-daemonset -o json | jq .status.numberAvailable)

# The return value of the script should contain a value and a map of labels
# value(int)|label1:value1,label2:value2
echo "$AVAILABLE"
exit 0