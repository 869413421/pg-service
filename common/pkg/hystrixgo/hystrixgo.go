package hystrixgo

import (
	"github.com/afex/hystrix-go/hystrix"
	"net"
	"net/http"
)

func HystrixBoot() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}
