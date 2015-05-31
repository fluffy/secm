
curl --data '{ "msg": "hello" }' --header "Content-Type:application/json" localhost:8081/v1/msg/55

curl --header "Accept:application/json" localhost:8081/v1/msg/1234-1


go run ms.go http://ks.fluffy.im:8080/
go run ms.go https://ks.fluffy.im/

