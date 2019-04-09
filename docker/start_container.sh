#!/bin/bash
echo "Starting container..."
if [ -z $(docker ps --filter=name='citrixitm_tf_dev_container' -q) ]
then
    if [ -z $(docker ps --filter=name='citrixitm_tf_dev_container' -a -q) ]
    then
        echo 'The citrixitm_tf_dev_container container doesn'\''t exist. Did you mean "make docker-run"?'
        exit 1
    else
        docker start -ai $1
    fi
else
    echo 'The citrixitm_tf_dev_container container is already running. Did you mean "make docker-exec-bash"?'
    exit 1
fi
