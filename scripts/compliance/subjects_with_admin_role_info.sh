#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check all subjects with admin role.
# INTERVAL 60

set -eux

# Retrieve all subjects with admin role
oc get rolebinding --all-namespaces -o json \
  | jq 'def getsubjects(data): data | "1|subject="+.subjects.kind+","+"name="+.subjects.name+","+"namespace="+.namespace; .items[] | select(.roleRef.name | test("^admin$")) | getsubjects({"subjects": .subjects[],"namespace": .metadata.namespace})' \
  | tr -d "\""

exit 0
