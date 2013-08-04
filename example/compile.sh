#!/bin/bash

#Make sure symlinks exist
./symlinker.sh

#Enable negative pattern matching !(prod*).go
shopt -s extglob

PRODFILES='./!(env_dev*).go'
DEVFILES='./!(env_prod*).go'

#Build dev binaries (we assume you are on a Mac) and prod binaries (which we assume is Linux)
#Change these to suit your needs
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o prod.forum.linux  ${PRODFILES} && ./gotogether prod.forum.linux "static templates" -q
GOOS=darwin GOARCH=amd64 go build -o dev.forum.osx  ${DEVFILES} && ./gotogether dev.forum.osx "static templates" -q

