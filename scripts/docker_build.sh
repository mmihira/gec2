#!/bin/sh
DIR="$( cd "$( dirname "$0" )" && pwd )"      # Dir of script location

docker build -t gec2:1.9 "$(dirname "$DIR")"

