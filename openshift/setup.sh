#!/bin/bash

oc new-project ose-healthchecker
oc new-app --docker-image=openshift-healthchecker:latest
oc new-build https://github.com/tobiasdenzler/openshift-healthchecker
oc new-app ose-healthchecker/openshift-healthchecker:latest
oc expose service openshift-healthchecker