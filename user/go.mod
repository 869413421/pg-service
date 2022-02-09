module github.com/869413421/pg-service/user

go 1.16

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/869413421/pg-service/common => ../common

require (
	github.com/869413421/pg-service/common v0.0.0-20220122101130-3dce972d21c7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.5.1
	github.com/thedevsaddam/govalidator v1.9.10
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/protobuf v1.27.1
)
