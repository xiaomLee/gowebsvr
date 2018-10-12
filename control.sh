#!/bin/bash
cwd=$(pwd)
workspace=$(cd $(dirname $0) && pwd -P)
cd $workspace

function start() {
    exec bin/gowebsvr -c ./conf/base.yaml
}

action=$1
case $action in
    "start" )
        start
        ;;
    * )
        echo "unknown command"
        exit 1
        ;;
esac
