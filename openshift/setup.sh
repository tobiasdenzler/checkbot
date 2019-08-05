#!/bin/bash

oc new-project ose-healthchecker
oc new-app --docker-image=openshift-healthchecker:latest
oc new-build https://github.com/tobiasdenzler/openshift-healthchecker
oc new-app ose-healthchecker/openshift-healthchecker:latest
oc expose service openshift-healthchecker

oc create sa healthchecker
oc adm policy add-role-to-user cluster-reader system:serviceaccount:healthchecker
oc patch dc/openshift-healthchecker --patch '{"spec":{"template":{"spec":{"serviceAccount":"healthchecker"}}}}'
oc patch cc/openshift-healthchecker --patch '{"spec":{"resources":{"requests":{"cpu":"200m"},{"memory":"200M"}},{"limits":{"cpu":"1"},{"memory":"1G"}}}}'

oc apply -f openshift/checker-daemonset.yaml
oc delete po -l name=checker-daemonset
