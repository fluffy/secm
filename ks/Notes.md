
start things up 
createMachine.sh test7
Remember to delete server at https://mycloud.rackspace.com

go run ks.go `docker-machine ip test7` 

Setting up GO
go get github.com/lib/pq

go run ks.go


curl http://localhost:8080/v1/key/100
curl --data "keyVal=fluffy&p=3" http://localhost:8080/v1/key



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

