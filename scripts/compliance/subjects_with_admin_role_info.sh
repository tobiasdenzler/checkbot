#!/bin/sh

# ACTIVE false
# TYPE Gauge
# HELP Check all subjects with admin role.
# INTERVAL 3600

set -eux

# Retrieve all subjects with admin role
oc get rolebinding --all-namespaces -o json \
  | jq 'def getsubjects(data): data | "1|subject="+.subjects.kind+","+"name="+.subjects.name+","+"namespace="+.namespace; .items[] | select(.roleRef.name | test("^admin$")) | getsubjects({"subjects": .subjects[],"namespace": .metadata.namespace})' \
  | tr -d "\""

exit 0
