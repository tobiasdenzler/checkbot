#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP get all openshift sccs where .allowPrivilegedContainer is set true
# INTERVAL 60

set -eux

# retrieve all sccs with .allowPrivilegedContainer is set true
oc get scc -o go-template='{{range .items}}{{if eq .allowPrivilegedContainer true}}{{$name := .metadata.name}}{{range $user := .users}}1|privileged_scc={{$name}},user={{$user}}{{"\n"}}{{end}}{{range $group := .groups}}1|privileged_scc={{$name}},group={{$group}}{{"\n"}}{{end}}{{end}}{{end}}'

echo "$(date -u +"%s")|lastrun=timestamp"

exit 0
