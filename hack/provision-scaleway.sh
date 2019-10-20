#!/bin/bash

export SUFFIX=$(head -c 16 /dev/urandom | shasum | cut -c1-8)
export NAME="inlets$SUFFIX"
export REGION="ams1" #(ams1 | par1)
export USERDATA=`pwd`/hack/userdata.sh
echo "Creating: $NAME"

rm ~/.scw-cache.db

scw --region=$REGION run --detach -u="cloud-init=@$USERDATA" --commercial-type=START1-XS --name=$NAME Ubuntu_Bionic

for i in {0..100};
do
  sleep 5
  status=$(scw --region=$REGION inspect -f {{.State}} server:"$NAME")
  if [ ! $? -eq 0 ];
  then
    echo "Unable to inspect server"
    rm ~/.scw-cache.db
    continue
  fi
  echo "Status: $status"

  if [ $status == "running" ];
  then
    echo IP: $(scw --region=$REGION inspect -f {{.PublicAddress.IP}} $NAME)
    break
  fi
done
