#!/bin/sh

# ACTIVE true
# TYPE Gauge
# HELP Simple check for testing.
# INTERVAL 10

set -eux

unknown

echo "42|label1=value1,label2=value2"
exit 0
