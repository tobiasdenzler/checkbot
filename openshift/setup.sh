#!/bin/bash

oc new-project ose-healthchecker
oc new-app --docker-image=openshift-healthchecker:latest