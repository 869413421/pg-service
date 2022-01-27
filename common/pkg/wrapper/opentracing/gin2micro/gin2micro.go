package gin2micro

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"math/rand"
	"net/http"
	"time"
)

const contextTracerKey = "Tracer-content"

var sf = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
func SetSamplingFrequency(n int) {
	sf = n
}

func TracerWrapper(c *gin.Context) {
	//1.从请求中的信息构建span作为链路起点
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
	tracer := opentracing.GlobalTracer()
	md := make(map[string]string)
	nsf := sf

	//2.判断请求头中是否包含上一个链路的信息，如果有作子SPAN，否则做根SPAN
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	if err == nil {
		sp = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		tracer = sp.Tracer()
		nsf = 100
	}
	defer sp.Finish()

	err = tracer.Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		return
	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	c.Set(contextTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > nsf {
		ext.SamplingPriority.Set(sp, 0)
	}
}

func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(contextTracerKey)
	if exist == false {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}
