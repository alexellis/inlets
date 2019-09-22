#!/usr/bin/env bash

declare -r here="$(dirname "$(pwd -P)/${BASH_SOURCE[0]}")"

SUFFIX=$(head -c 16 /dev/urandom | shasum | cut -c1-8)
INSTANCENAME="inlets-$SUFFIX"
SIZE="f1-micro"
ZONE="europe-west2-b"
USERDATA=${here}/userdata.sh

gcloud compute firewall-rules create allow-inlets \
    --allow tcp:80 \
    --source-ranges 0.0.0.0/0

echo "Creating: $INSTANCENAME"

output=$(gcloud compute instances create $INSTANCENAME \
    --machine-type=$SIZE \
    --metadata-from-file=startup-script=$USERDATA \
    --image-project="ubuntu-os-cloud" --image-family="ubuntu-1804-lts" \
    --zone=$ZONE \
    --format=json)

if [ $? -eq 0 ]
then

    name=$( echo $output | jq -r '.[0].name' )
    extip=$( echo $output | jq -r '.[0].networkInterfaces[0].accessConfigs[0].natIP' )

    echo "=============================="
    echo "Instance: $name has been created"
    echo "IP: $extip"
    echo "Login: gcloud compute ssh root@$name --zone=$ZONE"
    echo "=============================="
    echo "To destroy this instance run: gcloud compute instances delete $name --zone=$ZONE"
fi
