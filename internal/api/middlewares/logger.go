package middlewares

import (
	"backend/pkg/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func formatStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("\033[97;42m %d \033[0m", code)
	case code >= 300 && code < 400:
		return fmt.Sprintf("\033[97;46m %d \033[0m", code)
	case code >= 400 && code < 500:
		return fmt.Sprintf("\033[97;43m %d \033[0m", code)
	default:
		return fmt.Sprintf("\033[97;41m %d \033[0m", code)
	}
}

func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		requestID := GetRequestID(c)

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		if logger.Log == nil {
			return
		}

		latencyMS := float64(latency.Microseconds()) / 1000.0

		if !logger.IsProduction() {
			msg := fmt.Sprintf("%s %-6s %s  %.2fms", formatStatus(status), method, path, latencyMS)
			devLog := logger.Log.WithOptions(zap.WithCaller(false), zap.AddStacktrace(zapcore.DPanicLevel))

			var fields []zap.Field
			if len(c.Errors) > 0 {
				fields = append(fields, zap.String("error", c.Errors.String()))
			}
			if status >= 400 {
				fields = append(fields, zap.String("req_id", requestID), zap.String("client_ip", clientIP))
			}

			switch {
			case status >= 500:
				devLog.Error(msg, fields...)
			case status >= 400:
				devLog.Warn(msg, fields...)
			default:
				devLog.Info(msg, fields...)
			}
			return
		}

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Float64("latency_ms", latencyMS),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		msg := "HTTP Request"
		switch {
		case status >= 500:
			logger.Log.Error(msg, fields...)
		case status >= 400:
			logger.Log.Warn(msg, fields...)
		default:
			logger.Log.Info(msg, fields...)
		}
	}
}
