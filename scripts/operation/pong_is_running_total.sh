#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check if all pods from Daemonset are running.
# INTERVAL 30

set -eux

AVAILABLE=$(oc get ds/pong-daemonset -n checkbot -o json | jq .status.numberAvailable)

# The return value of the script should contain a value and a map of labels
# value|label1=value1,label2=value2
echo "$AVAILABLE"
exit 0
