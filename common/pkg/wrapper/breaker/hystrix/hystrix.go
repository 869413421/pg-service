package hystrix

import (
	"context"
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/micro/go-micro/v2/client"
	"net"
	"net/http"
	"time"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 初始化retrier，每隔1000ms重试一次，总共重试3次
		r := retrier.New(retrier.ConstantBackoff(3, 1000 * time.Millisecond), nil)
		// retrier 工作模式和 hystrix 类似，在 Run 方法中将待执行的业务逻辑封装到匿名函数传入即可
		err := r.Run(func() error {
			return c.Client.Call(ctx, req, rsp, opts...)
		})
		return err
	}, func(err error) error {
		// 你可以在这里自定义更复杂的服务降级逻辑作为服务熔断的兜底
		logger.Danger("hystrix fallback error: %v", err)
		return err
	})
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}

func Configure(names []string) {
	// 1.初始化熔断设置
	config := hystrix.CommandConfig{
		//用于设置超时时间，超过该时间没有返回响应，意味着请求失败
		Timeout: 3000,
		//用于设置同一类型请求的最大并发量，达到最大并发量后，接下来的请求会被拒绝
		MaxConcurrentRequests: 100,
		//用于设置指定时间窗口内让断路器跳闸（开启）的最小请求数；
		RequestVolumeThreshold: 100,
		//断路器跳闸后，在此时间段内，新的请求都会被拒绝；
		SleepWindow: 3000,
		//请求失败百分比，如果超过这个百分比，则断路器跳闸
		ErrorPercentThreshold: 50,
	}

	//2.根据服务名加载默认配置
	configs := make(map[string]hystrix.CommandConfig)
	for _, name := range names {
		configs[name] = config
	}
	hystrix.Configure(configs)

	//3.hystrix监听数据
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
}
