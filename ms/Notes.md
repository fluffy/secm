
curl --data '{ "msg": "hello" }' --header "Content-Type:application/json" localhost:8081/v1/msg/55

curl --header "Accept:application/json" localhost:8081/v1/msg/1234-1


The / on end of ksURL on next line is required
go run ms.go http://ks.fluffy.im:8080/ ks.fluffy.im
go run ms.go https://ks.fluffy.im/ ks.fluffy.im
