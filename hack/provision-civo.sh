#!/bin/bash

export SUFFIX=$(head -c 16 /dev/urandom | shasum | cut -c1-8)
export INSTANCENAME="inlets$SUFFIX"
export SIZE="g2.xsmall"
export TEMPLATE_ID="811a8dfb-8202-49ad-b1ef-1e6320b20497" # ubuntu-16.04
export REGION="lon1"
export FIELDS="ID,Name,PublicIPv4"
export SCRIPT=`pwd`/hack/userdata.sh

echo "Creating: $INSTANCENAME"

instanceId="$(civo instance create --name=$INSTANCENAME \
               --size=$SIZE \
               --user=root \
               --template-id=$TEMPLATE \
               --script=$SCRIPT \
               --wait \
               --quiet \
               )"
               
if [ $? -eq 0 ];
then

instanceInfo="$(civo instance show $instanceId)"
publicIp="$(echo $instanceInfo | grep "Public IP" | cut -d"=" -f2 | cut -c3-)"
name="$(echo $instanceInfo | grep "Hostname" | cut -d":" -f2 | cut -c2-)"

echo "=============================="
echo "Instance: $name has been created"
echo "IP: $publicIp"
echo "Login: ssh root@$publicIp"
echo "=============================="
echo "To destroy this instance run: civo instance rm $instanceId"

fi
