#!/bin/sh
PROJECTS=$(oc get project --no-headers | wc -l)
QUOTAS=$(oc get quota --all-namespaces --no-headers | wc -l)

if [ $PROJECTS -eq $QUOTAS ]
then
  echo "ok"
  exit 0
fi
echo "found $PROJECTS project and $QUOTAS quotas"
exit 1