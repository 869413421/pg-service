module github.com/869413421/pg-service/user-api

go 1.13

replace github.com/869413421/pg-service/common => ../common

require (
	github.com/869413421/pg-service/common v0.0.0-20220125095543-935236a3185e
	github.com/869413421/pg-service/user v0.0.0-20220127033843-d721ccd1897a
	github.com/gin-gonic/gin v1.7.7
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/breaker/hystrix/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.1.0
)
