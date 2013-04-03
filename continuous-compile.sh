#!/bin/bash

# This script keeps watch on the current project and compiles it continuously as you change files.
# If there are multiple projects with the same final directory name (e.g., /proj/rad and /lib/monster/rad), 
# this will kill any other similarly-named running projects' binaries, potentially leading to havoc.
 
# To run, execute the following:
# /usr/local/bin/fswatch ./ ./continuous-compile.sh

echo "Re-compiling"

#Count the number of times the compiled app is running
numproc=`ps aux | grep ${PWD##*/}-main.osx | grep -v grep | wc -l`

if [[ numproc > 0 ]]
	then
	killall ${PWD##*/}-main.osx 2> /dev/null
fi

go build -o /tmp/${PWD##*/}-main.osx *.go && /tmp/${PWD##*/}-main.osx &

