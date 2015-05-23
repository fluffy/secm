#!/bin/bash 

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
docker-machine create --driver rackspace --rackspace-flavor-id 2   test3

# TODO - remove public port 
docker `docker-machine config test3` run -p 5432:5432 --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres 




