#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check all subjects with cluster-admin role.
# INTERVAL 60

set -eux

# Retrieve all subjects with cluster-admin role
SUBJECTS=$(oc get clusterrolebinding -o json | jq '.items[] | select(.metadata.name |  startswith("cluster-admin")) | .subjects[] | "subject="+.kind+","+"name="+.name')

for subject in $SUBJECTS
do
    # Return the subjects with cluster-admin role, tr will strip quotes
    # The return value of the script should contain a value and a map of labels
    # value|label1=value1,label2=value2
    echo "1|$subject" | tr -d "\""
done

exit 0
