#!/bin/bash

oc new-project ose-healthchecker
oc new-app --docker-image=openshift-healthchecker:latest
oc new-build https://github.com/tobiasdenzler/openshift-healthchecker
oc new-app ose-healthchecker/openshift-healthchecker:latest
oc expose service openshift-healthchecker

oc create sa healthchecker
oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:ose-healthchecker:healthchecker
oc patch dc/openshift-healthchecker --patch '{"spec":{"template":{"spec":{"serviceAccount":"healthchecker"}}}}'
oc patch dc/openshift-healthchecker --patch '{"spec":{"resources":{"requests":{"cpu":"200m"},{"memory":"200M"}},{"limits":{"cpu":"1"},{"memory":"1G"}}}}'

oc apply -f openshift/setup

oc delete po -l name=checker-daemonset
