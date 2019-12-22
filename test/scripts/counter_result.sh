#!/bin/sh

# ACTIVE true
# TYPE Counter
# HELP Simple check for testing.
# INTERVAL 10

set -eux

echo "1|label1=value1,label2=value2"
exit 0
