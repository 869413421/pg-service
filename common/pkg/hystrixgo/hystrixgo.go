package hystrixgo

import (
    "github.com/afex/hystrix-go/hystrix"
	"net"
	"net/http"
)

func StartHystrixClient() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
}
