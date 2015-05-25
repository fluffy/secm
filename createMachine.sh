#!/bin/bash 
set -e

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


# flavors can be found at http://docs.rackspace.com/cas/api/v1.0/autoscale-devguide/content/server-flavors.html
docker-machine create --driver rackspace --rackspace-flavor-id 2 "$MAC_NAME"

#build the docker images 
docker-machine ssh $MAC_NAME git clone https://github.com/fluffy/secm.git
docker-machine ssh $MAC_NAME "docker build -t fluffy/ks:v1 /root/secm/ks"
docker-machine ssh $MAC_NAME "docker build -t fluffy/ws:v1 /root/secm/ws"

# start the varios machins
# TODO - remove public port 
docker `docker-machine config $MAC_NAME` run -p 5432:5432 --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres
# TODO - remove public port
docker `docker-machine config $MAC_NAME` run -p 8080:8080  --name="my-ks" --link my-postgres:db -d fluffy/ks:v1
# TODO remove the ssh port 22
docker `docker-machine config $MAC_NAME` run -p 80:80 -p 443:443 -p 8022:22  -v /root/data:/data --name="my-ws" --link my-ks:ks -d fluffy/ws:v1

docker-machine ssh $MAC_NAME mkdir /root/data
cat site.crt | docker-machine ssh $MAC_NAME "cat > /root/data/site.crt"
cat site-chain.crt | docker-machine ssh $MAC_NAME "cat > /root/data/site-chain.crt"
cat site.key | docker-machine ssh $MAC_NAME "cat > /root/data/site.key"
cat site.conf | \
    sed -e "s/SERVER_NAME/$SECM_KS_NAME/g" | \
    sed -e "s/OIDC_Client_ID/$OIDC_Client_ID/g" | \
    sed -e "s/OIDC_Client_Secret/$OIDC_Client_Secret/g" | \
    sed -e "s/OIDC_Crypto_Passphrase/$OIDC_Crypto_Passphrase/g" | \
    docker-machine ssh $MAC_NAME "cat > /root/data/site.conf"

echo
echo Machine IP is `docker-machine ip $MAC_NAME`
echo Remember to delete server at https://mycloud.rackspace.com


