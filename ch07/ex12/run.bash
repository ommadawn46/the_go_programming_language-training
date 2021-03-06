#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

trap "exit" INT TERM ERR
trap "kill 0" EXIT

go run stockServer.go &
sleep 3s

curl 'http://localhost:8000/list'
