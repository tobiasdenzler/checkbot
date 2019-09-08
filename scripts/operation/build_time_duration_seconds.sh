#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Measure the time a build will run.
# INTERVAL 900

set -eux

start=$(date +%s)
oc start-build -w pong -n checkbot
stop=$(date +%s)
elapsed="$(($stop-$start))"

# return the number of seconds the build was running
echo "$elapsed"
exit 0
