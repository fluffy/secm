
start things up 
createMachine.sh test7
Remember to delete server at https://mycloud.rackspace.com

go run ks.go `docker-machine ip test7` 

Setting up GO
go get github.com/lib/pq
go get github.com/gorilla/mux

go run ks.go


curl http://localhost:8080/v1/key/100
curl --data "keyVal=fluffy&p=3" http://localhost:8080/v1/key


curl --data "keyVal=__2__" http://ks.fluffy.im:8080/v1/key


curl http://ks.fluffy.im:8080/v1/key/861802244120027882



Info for docker on rackspace at
https://developer.rackspace.com/blog/using-docker-machine-to-deploy-your-docker-containers-on-rackspace/

Install Docker Machine (see)
https://github.com/docker/machine/releases


Set up variables for rackspace 
OS_USERNAME, OS_API_KEY and OS_REGION_NAME

Set up variable for database 
SECM_DB_SECRET


Make a new docker machine 
docker-machine create --driver rackspace test3

# flavor  can be found at http://docs.rackspace.com/cas/api/v1.0/autoscale-devguide/content/server-flavors.html
docker-machine create --driver rackspace --rackspace-flavor-id 2   test3


Set up environ variables returned by
docker-machine env test3

docker-machine ssh test3 "mkdir /data"

docker `docker-machine config test3` run --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres 

#this one makes the database port public
docker `docker-machine config test3` run -p 5432:5432 --name my-postgres -e POSTGRES_PASSWORD=$SECM_DB_SECRET -d postgres 

docker `docker-machine config test3` ps
docker-machine ip test3


# note default database called postgres 
docker `docker-machine config test3` run -it --link my-postgres:postgres --rm postgres sh -c 'PGPASSWORD='$SECM_DB_SECRET' exec psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres -w '

# or my local mac
setenv PGPASSWORD  $SECM_DB_SECRET
/Library/PostgreSQL/9.1/bin/psql -h `docker-machine ip test7` -p 5432 -U postgres


/Library/PostgreSQL/9.1/bin/pg_dump -h `docker-machine ip test7` -p 5432 -U postgres -w postgres



#CREATE USER dbUser UNENCRYPTED PASSWORD 'test' ;
#CREATE DATABASE myDb OWNER dbUser;

#DROP TABLE keys;
#DROP TABLE keyUsers;
#DROP TABLE keyAdmins;


# should add dates to keys so can expire them
# should encrypt kVal and add epoc so can rotate the keys used so that DB and DB backup are encrypted

echo \
    "CREATE TABLE keys ( kID BIGINT NOT NULL, kVal bytea NOT NULL ,  oID BIGINT NOT NULL, PRIMARY KEY( kID ) );" \
    "CREATE TABLE keyUsers ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );" \
    "CREATE TABLE keyAdmins ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );" \
    | docker `docker-machine config test3` run -it --link my-postgres:postgres --rm postgres sh -c 'PGPASSWORD='$SECM_DB_SECRET' exec psql -h "$POSTGRES_PORT_5432_TCP_ADDR" -p "$POSTGRES_PORT_5432_TCP_PORT" -U postgres -w -f- '

CREATE TABLE keys ( kID BIGINT NOT NULL, kVal bytea NOT NULL ,  oID BIGINT NOT NULL, PRIMARY KEY( kID ) );
CREATE TABLE keyUsers ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );
CREATE TABLE keyAdmins ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );

INSERT INTO keys (kID, kVal, oID) VALUES (100,'a',1);
INSERT INTO keyUsers (kID,uID) VALUES (100,1);
INSERT INTO keyAdmins (kID,uID) VALUES (100,1);

INSERT INTO keys (kID, kVal, oID) VALUES (101,'b',1);
INSERT INTO keyUsers (kID,uID) VALUES (101,1);
INSERT INTO keyAdmins (kID,uID) VALUES (101,1);

INSERT INTO keys (kID, kVal, oID) VALUES (102,'c',2);
INSERT INTO keyUsers (kID,uID) VALUES (102,2);
INSERT INTO keyAdmins (kID,uID) VALUES (102,2);

INSERT INTO keyUsers (kID,uID) VALUES (100,2);



SELECT keys.kVal FROM keys;

SELECT kVal FROM keys WHERE kID=100;


SELECT * FROM keyUsers WHERE kID=100 AND uID=2;


SELECT * FROM keys, keyUsers WHERE keys.kID = keyUsers.kID;


SELECT keys.kVal FROM keys, keyUsers WHERE keys.kID = 100 AND keyUsers.kID = 100 AND keyUsers.uID = 2 ;

SELECT keys.kVal  FROM keys JOIN keyUsers ON  keys.kID = keyUsers.kID WHERE keyUsers.uID = 2 AND keyUsers.kID = 102 ;

# if user 2 is a user of key 102, return the key value 
SELECT keys.kVal  FROM keyUsers JOIN keys ON  keys.kID = keyUsers.kID WHERE keyUsers.uID = 2 AND keyUsers.kID = 102 ;

# if uid 2 is an admin for key 102, then add 3 as a user of this key 
INSERT INTO keyUsers (kID,uID) SELECT kID,3  FROM keyAdmins WHERE keyAdmins.kID = 102 AND keyAdmins.uID = 2;

# if uid 1 is an owner for key 101, then add 3 as a admin of this key 
INSERT INTO keyAdmins (kID,uID) SELECT kID,3  FROM keys WHERE keys.kID = 101 AND keys.oID = 1;


# if uid 3 is a user of key 100, then return who the owner or that key is 
SELECT keys.oID FROM keys JOIN  keyUsers ON keys.kID = keyUsers.kID WHERE keys.kID = 100 AND keyUsers.uID = 3  ;

# if uid 2 is a user of key 100, then return all the admins of that key 
SELECT keyAdmins.uID FROM keyAdmins JOIN keyUsers ON keyAdmins.kID = keyUsers.kID WHERE keyUsers.kID = 100 AND keyUsers.uID = 2  ;

# if uid 2 is a uer of key 100, then return all the users of that key 
SELECT users.uID FROM keyUsers AS users JOIN keyUsers AS perm ON users.kID = perm.kID WHERE perm.kID = 100 AND perm.uID = 2  ;



building mod_auth_openidc

sudo apt-get install autotools-dev
sudo apt-get install autoconf
sudo apt-get install apache-dev
sudo apt-get install apache2-dev
sudo apt-get install curl
sudo apt-get install jansson
sudo apt-get install pcre3
sudo apt-get install libcurl
sudo apt-get install curl libc6 libcurl3 zlib1g
sudo apt-get install libjansson-dev
sudo apt-get install libpcre3 libpcre3-dev
sudo apt-get install pkgconfig
sudo apt-get install pkg-config
sudo apt-get install redis
sudo apt-get install hiredis
sudo apt-get install libcurl4-openssl-dev

./configure

make
sudo make install 



curl --data "keyVal=x2x" http://192.237.200.55:8080/v1/key

curl http://192.237.200.55:8080/v1/key/7699388317967658784



OIDC_CLAIM_sub: xxx
OIDC_CLAIM_profile: https://plus.google.com/xxx
OIDC_CLAIM_gender: male
OIDC_CLAIM_family_name: xx
OIDC_CLAIM_name: xx xx
OIDC_CLAIM_given_name: xx
OIDC_CLAIM_picture: https://lh4.googleusercontent.com/xxx/photo.jpg
OIDC_CLAIM_email: xxx@gmail.com
OIDC_CLAIM_email_verified: 1
OIDC_CLAIM_iss: https://accounts.google.com
OIDC_CLAIM_nonce: xxx
OIDC_CLAIM_exp: 1432588106
OIDC_CLAIM_iat: 1432584506
OIDC_CLAIM_at_hash: xx
OIDC_CLAIM_azp: xxx
OIDC_CLAIM_aud: xxx
OIDC_access_token: xxx
OIDC_access_token_expires: 1432588107
X-Forwarded-For: 10.1.3.233
X-Forwarded-Host: ks.fluffy.im
X-Forwarded-Server: ks.fluffy.im


curl --header "OIDC_CLAIM_email: xxx@gmail" http://localhost:8080/



curl --data "__5__" --header "OIDC_CLAIM_email: xxx@gmail" --header
"OIDC_CLAIM_email_verified: 1" http://localhost:8080/v1/key

curl --header "OIDC_CLAIM_email: xxx@gmail" --header "OIDC_CLAIM_email_verified:
1" http://localhost:8080/v1/key/2991017719384258649



