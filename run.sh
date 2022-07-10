#!/bin/bash
#
rm logs
botname="CheckNewListings"

echo "botname = $botname"

newprogram="./$botname"

PID=`ps -aef | grep " $newprogram$" | grep -v grep | awk '{print $2}'`
echo "pid = $PID"

if [ "$PID" != ""  ]; then 
	echo there is already a program run stop.sh
	exit 3
fi

# git fetch
# if [ $? -eq 0 ]; then
#   git checkout $branch
#   if [ $? -eq 0 ]; then
#     git pull
#     if [ $? -eq 0 ]; then

      go build -o "$botname"
#     fi
#   fi
# fi
"./$botname" &

sleep 1

disown

echo OK


