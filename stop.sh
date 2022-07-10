#!/bin/bash
#

pid=`ps -aef | grep ./CheckNewListing | grep -v grep | awk '{print $2}'`
echo "pid = $pid"

if [ "$pid" = "" ]; then
  echo no such program
else
  echo OK
    kill $pid
fi
