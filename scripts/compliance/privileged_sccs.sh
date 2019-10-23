#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP get all openshift sccs where .allowPrivilegedContainer is set true
# INTERVAL 60

set -eux

# retrieve all sccs with .allowPrivilegedContainer is set true
oc get scc -o go-template='{{range .items}}{{if eq .allowPrivilegedContainer true}}1|privileged_scc_name={{.metadata.name}}{{"\n"}}{{end}}{{end}}'

exit 0
