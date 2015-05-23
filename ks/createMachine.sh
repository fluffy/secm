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


