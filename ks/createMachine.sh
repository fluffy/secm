#!/bin/bash 

if [ -z "$1" ]; then
    echo usage: $0 MAC_NAME;
    exit;
fi
MAC_NAME="$1"

if [ -z "$SECM_DB_SECRET" ]; then
    echo need to set SECM_DB_SECRET;
    exit;
fi

if [ -z "$OS_USERNAME" ]; then
    echo need to set OS_USERNAME with rackspace username;
    exit;
fi

if [ -z "$OS_API_KEY" ]; then
    echo need to set OS_API_KEY  with rackspace API_KEY;
    exit;
fi

if [ -z "$OS_REGION_NAME" ]; then
    echo need to set OS_REGION_NAME with rackspace REGION_NAME;
    exit;
fi


# flavor  can be found at http://docs.rackspace.com/cas/api/v1.0/autoscale-devguide/content/server-flavors.html
docker-machine create --driver rackspace --rackspace-flavor-id 2 "$MAC_NAME"

# TODO - remove public port 
docker `docker-machine config $MAC_NAME` run -p 5432:5432 --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres 

echo Remember to delete server at https://mycloud.rackspace.com






