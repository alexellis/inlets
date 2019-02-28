#!/bin/sh

export eTAG="latest"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

Version=$(git describe --tags --dirty)
GitCommit=$(git rev-parse HEAD)


echo Building alexellis2/inlets:$eTAG

docker build --build-arg VERSION=$Version --build-arg GIT_COMMIT=$GitCommit -t alexellis2/inlets:$eTAG . && \
 docker create --name inlets alexellis2/inlets:$eTAG && \
 docker cp inlets:/root/inlets . && \
 docker rm -f inlets
