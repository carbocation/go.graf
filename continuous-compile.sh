#!/bin/bash

#Enable negative pattern matching !(prod*).go
shopt -s extglob

# This script keeps watch on the current project and compiles it continuously as you change files.
# If there are multiple projects with the same final directory name (e.g., /proj/rad and /lib/monster/rad), 
# this will kill any other similarly-named running projects' binaries, potentially leading to havoc.
 
# To run, execute the following:
# /usr/local/bin/fswatch ./ ./continuous-compile.sh

#NOT the production environment
NOTENV='./!(env_prod*).go'
echo 'Re-compiling with these options: '"$@"

#Count the number of times the compiled app is running
numproc=`ps aux | grep ${PWD##*/}-main.osx | grep -v grep | wc -l`

if [[ numproc > 0 ]]
	then
	killall ${PWD##*/}-main.osx 2> /dev/null
fi

# Build the binary, then
# add in the static components with NRSC, then
# launch the binary in the background
GOMAXPROCS=4 go build "$@" -o /tmp/${PWD##*/}-main.osx ${NOTENV} && ./nrsc-script /tmp/${PWD##*/}-main.osx "static templates" -q && /tmp/${PWD##*/}-main.osx &
