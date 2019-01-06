#!/bin/bash

export SUFFIX=$(head -c 16 /dev/urandom | shasum | cut -c1-8)
export DROPLETNAME="inlets$SUFFIX"
export SIZE="s-1vcpu-1gb"
export IMAGE="ubuntu-16-04-x64"
export REGION="lon1"
export FIELDS="ID,Name,PublicIPv4"
export USERDATA=`pwd`/hack/userdata.sh

echo "Creating: $DROPLETNAME"

dropletInfo="$(doctl compute droplet create $DROPLETNAME \
               --size $SIZE \
               --image $IMAGE \
               --region $REGION \
               --user-data-file $USERDATA \
               --format "$FIELDS" \
               --no-header \
               --wait \
               )"
               
if [ $? -eq 0 ];
then

readfields=$(sed 's/,/ /g' <<<$FIELDS)
read -r $readfields <<<"$dropletInfo"

echo "=============================="
echo "Droplet: $Name has been created"
echo "IP: $PublicIPv4"
echo "Login: ssh root@$PublicIPv4"
echo "=============================="
echo "To destroy this droplet run: doctl compute droplet delete -f $ID"

fi