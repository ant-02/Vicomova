package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 忽略的路径（不记录指标）
	ignorePaths = map[string]bool{
		"/metrics": true,
		"/health":  true,
	}
)

// PrometheusMiddleware Prometheus 指标中间件
func PrometheusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		if ignorePaths[path] {
			c.Next(ctx)
			return
		}

		start := time.Now()
		method := string(c.Request.Method())

		c.Next(ctx)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response.StatusCode())
		pathNormalized := normalizePath(path)

		HTTPRequestsTotal.WithLabelValues(method, pathNormalized, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, pathNormalized).Observe(duration)
	}
}

// normalizePath 标准化路径
func normalizePath(path string) string {
	// TODO: 进一步处理路径参数化
	return path
}

// MetricsHandler 返回 Prometheus metrics 的 HTTP handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
