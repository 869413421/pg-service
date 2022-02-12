package prometheus

import (
	"github.com/869413421/pg-service/common/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// PrometheusBoot 启动 HTTP 服务监听客户端数据采集
func PrometheusBoot() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9092", nil)
		if err != nil {
			logger.Danger("ListenAndServe: ", err)
		}
	}()
}
