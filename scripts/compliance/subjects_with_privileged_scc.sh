#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP get all users and groups in openshift sccs where .allowPrivilegedContainer is set true
# INTERVAL 60

set -eux

# retrieve all users and groups in sccs with .allowPrivilegedContainer is set true
oc get scc -o go-template='{{range .items}}{{if eq .allowPrivilegedContainer true}}{{$name := .metadata.name}}{{range $user := .users}}1|privileged_scc={{$name}},subject={{$user}},type=user{{"\n"}}{{end}}{{range $group := .groups}}1|privileged_scc={{$name}},subject={{$group}},type=group{{"\n"}}{{end}}{{end}}{{end}}'

exit 0
