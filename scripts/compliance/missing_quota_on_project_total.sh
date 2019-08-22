#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 60

set -eux

PROJECTS=$(oc get project --no-headers | wc -l)
QUOTAS=$(oc get quota --all-namespaces --no-headers | wc -l)

DIFF="$(($PROJECTS-$QUOTAS))"

# The return value of the script should contain a value and a map of labels
# value(int)|label1=value1,label2=value2
echo "$DIFF"
exit 0