#!/bin/bash

export SUFFIX=$(head -c 16 /dev/urandom | shasum | cut -c1-8)
export NAME="inlets$SUFFIX"
export USERDATA=$(cat `pwd`/hack/userdata.sh)

export VULTR_REGION=8 #London
export VULTR_OS=270 #ubuntu 18.04
export VULTR_PLAN=201 #1024 MB RAM,25 GB SSD,1.00 TB BW

echo "Creating: $NAME"

VULTR_SCRIPT_ID=$(vultr-cli script create --name $NAME --script "$USERDATA" | rev | cut -f1 -d ' ' | rev)
VULTR_SERVER_ID=$(vultr-cli server create --os $VULTR_OS --plan $VULTR_PLAN --region $VULTR_REGION --host $NAME --script-id $VULTR_SCRIPT_ID | rev | cut -f1 -d ' ' | rev)

for i in {0..120};
do
  sleep 5
  status=$(vultr-cli server info $VULTR_SERVER_ID | grep 'Server State' | rev | cut -f1 | rev)
  if [ ! $? -eq 0 ];
  then
    echo "Unable to inspect server"
    continue
  fi
  echo "Status: $status"

  if [ $status == "ok" ];
  then
    IP=$(vultr-cli server info $VULTR_SERVER_ID | grep 'Main IP' | rev | cut -f1 | rev)

    echo "=============================="
    echo "Server: $VULTR_SERVER_ID has been created"
    echo "IP: $IP"
    echo "Login: ssh root@$IP"
    echo "=============================="
    echo "To destroy this instance run: vultr-cli server delete $VULTR_SERVER_ID"

    break
  fi
done