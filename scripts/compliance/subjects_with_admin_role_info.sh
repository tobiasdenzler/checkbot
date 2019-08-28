#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Check all subjects with admin role.
# INTERVAL 60

set -eux

# Retrieve all subjects with admin role
SUBJECTS=$(oc get rolebinding --all-namespaces -o json | jq 'def getsubjects(data): data | "subject="+.subjects.kind+","+"name="+.subjects.name+","+"namespace="+.namespace; .items[] | select(.roleRef.name | test("admin")) | getsubjects({"subjects": .subjects[],"namespace": .metadata.namespace})')

for subject in $SUBJECTS
do
    # Return the subjects with admin role, tr will strip quotes
    echo "1|$subject" | tr -d "\""
done

exit 0
