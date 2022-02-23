module github.com/869413421/pg-service/common

go 1.15

require (
	github.com/869413421/pg-service/user v0.0.0-20220127033843-d721ccd1897a
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/eapache/go-resiliency v1.1.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-acme/lego/v3 v3.4.0
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.12
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/broker/rabbitmq/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.1.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)
