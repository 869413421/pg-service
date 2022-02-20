package hystrixgo

import (
	"github.com/afex/hystrix-go/hystrix"
	"net"
	"net/http"
)

func StartHystrixClient() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		http.ListenAndServe(net.JoinHostPort("", "2020"), hystrixStreamHandler)
	}()
}
