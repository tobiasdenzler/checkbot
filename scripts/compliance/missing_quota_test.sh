#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check if all projects have quotas defined.
# INTERVAL 10

set -eux

PROJECTS=6
QUOTAS=4
DIFF="$(($PROJECTS-$QUOTAS))"

# The return value of the script should contain a value and a map of labels
# value(int)|label1:value1,label2:value2
echo "$DIFF|label1:value1,label2:value2"
#echo "$DIFF|label1:value1,label2:value2"
exit 0
