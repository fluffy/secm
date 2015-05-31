#!/bin/bash 
set -e

echo -n Starting: ; date 


if [ -z "$1" ]; then
    echo usage: $0 MAC_NAME;
    exit;
fi
MAC_NAME="$1"

if [ -z "$SECM_DB_SECRET" ]; then
    echo need to set SECM_DB_SECRET with desired password for database;
    exit;
fi

if [ -z "$OS_USERNAME" ]; then
    echo need to set OS_USERNAME with rackspace username;
    exit;
fi

if [ -z "$OS_API_KEY" ]; then
    echo need to set OS_API_KEY with rackspace API_KEY;
    exit;
fi

if [ -z "$OS_REGION_NAME" ]; then
    echo need to set OS_REGION_NAME with rackspace REGION_NAME;
    exit;
fi

# TODO - check $RS_ACCOUNT_NUMBER and $RS_DOMAIN_ID exist

# TODO - check that the files with cert, key, and chain exist 

# flavors can be found at http://docs.rackspace.com/cas/api/v1.0/autoscale-devguide/content/server-flavors.html
docker-machine create --driver rackspace --rackspace-flavor-id 2 "$MAC_NAME"

echo Machine IP is `docker-machine ip $MAC_NAME`

#set the IP to point at it
export TOKEN=` curl -D - -H "X-Auth-Key: $OS_API_KEY" -H "X-Auth-User: $OS_USERNAME" https://auth.api.rackspacecloud.com/v1.0 | grep "X-Auth-Token\:" | awk ' { print $2 } ' `
export IP=`docker-machine ip $MAC_NAME`
curl -X PUT -H X-Auth-Token:\ $TOKEN -H Content-Type:\ application/json https://dns.api.rackspacecloud.com/v1.0/$RS_ACCOUNT_NUMBER/domains/$RS_DOMAIN_ID/records/$RS_RECORD --data ' { "data": "'$IP'" } '

#build the docker images 
docker-machine ssh $MAC_NAME git clone https://github.com/fluffy/secm.git
docker-machine ssh $MAC_NAME "docker build -t fluffy/ks:v1 /root/secm/ks"
docker-machine ssh $MAC_NAME "docker build -t fluffy/ws:v1 /root/secm/ws"

# set up the site specific data 
docker-machine ssh $MAC_NAME mkdir /root/data
cat ws/site.crt | docker-machine ssh $MAC_NAME "cat > /root/data/site.crt"
cat ws/site-chain.crt | docker-machine ssh $MAC_NAME "cat > /root/data/site-chain.crt"
cat ws/site.key | docker-machine ssh $MAC_NAME "cat > /root/data/site.key"
cat ws/site.conf | \
    sed -e "s/SERVER_NAME/$SECM_KS_NAME/g" | \
    sed -e "s/OIDC_Client_ID/$OIDC_Client_ID/g" | \
    sed -e "s/OIDC_Client_Secret/$OIDC_Client_Secret/g" | \
    sed -e "s/OIDC_Crypto_Passphrase/$OIDC_Crypto_Passphrase/g" | \
    docker-machine ssh $MAC_NAME "cat > /root/data/site.conf"

# start the varios machines
# TODO - remove public ports
docker `docker-machine config $MAC_NAME` run -p 5432:5432   --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres
docker `docker-machine config $MAC_NAME` run -p 27017:27017 --name my-mongo -d mongo
docker `docker-machine config $MAC_NAME` run -p 6379:6379   --name my-redis -d redis
docker `docker-machine config $MAC_NAME` run -p 9042:9042   --name my-cassandra -d cassandra
sleep 10

# TODO - remove public port
docker `docker-machine config $MAC_NAME` run -p 8080:8080  --name="my-ks" --link my-postgres:db -d fluffy/ks:v1
sleep 5
docker `docker-machine config $MAC_NAME` run -p 80:80 -p 443:443 -v /root/data:/data --name="my-ws" --link my-ks:ks -d fluffy/ws:v1
sleep 5

echo
echo Machine IP is `docker-machine ip $MAC_NAME`
echo Remember to delete server at https://mycloud.rackspace.com


echo -n Finished ; date

