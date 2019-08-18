#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check all users with cluster-admin role.
# INTERVAL 20

set -eux

# Retrieve all users with cluster-admin
USERS=$(oc get clusterrolebinding -o json | jq '.items[] | select(.metadata.name |  startswith("cluster-admin")) | if .subjects[].kind == "User" then .subjects[].name else empty end')

for user in $USERS
do
    # Return the users with cluster-admin role, tr will strip quotes
    echo "1|user:$user" | tr -d "\""
done

exit 0