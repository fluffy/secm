#!/bin/bash 

if [ -z "$SECM_DB_SECRET" ]; then
    echo need to set SECM_DB_SECRET;
    exit;
fi


