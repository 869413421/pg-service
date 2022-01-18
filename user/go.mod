module github.com/869413421/pg-service/user

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/869413421/pg-service/common => ../common

require (
	github.com/869413421/pg-service/common v0.0.0-20220116122049-a771def830a7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro/v2 v2.9.1
	github.com/thedevsaddam/govalidator v1.9.10
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	google.golang.org/protobuf v1.27.1
)
