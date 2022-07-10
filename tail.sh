#!/bin/bash
#
if [ "$1" = "" ]; then 
	exit 1
fi

dir=$1
basedir=$PWD
cd $dir
if [ $? -ne 0 ]; then 
	exit 1
fi

tail -fn 350 logs


