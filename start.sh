#!/bin/sh

echo "* * * * * echo '123'" > cronfile

# Load the crontab file
crontab cronfile

crond -f