package logger

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxLoggedBody = 4096 // จำกัดความยาวที่เก็บใน log กันไฟล์โตเกิน

// moduleFromRoute แปลง route เป็นชื่อโมดูล เช่น "/books/:id" -> "books"
func moduleFromRoute(route string) string {
	route = strings.Trim(route, "/")
	if route == "" {
		return "api"
	}
	return strings.SplitN(route, "/", 2)[0]
}

// bodyLogWriter ดัก response body ที่ framework จะส่งออก เพื่อเก็บสำเนาไปลง log
type bodyLogWriter struct {
	gin.ResponseWriter
	buffer bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.buffer.Write(b) // เก็บสำเนาไว้ใน buffer
	return w.ResponseWriter.Write(b)
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLoggedBody {
		return s[:maxLoggedBody] + "..."
	}
	return s
}

// AccessLog เขียน log ทั้ง request + response ทุกสถานะ
// 2xx → info, 4xx → warn, 5xx → error
func AccessLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		// เก็บ request body (จำกัดความยาว)
		var requestBody string
		if context.Request != nil && context.Request.Body != nil && context.Request.ContentLength != 0 {
			raw, _ := io.ReadAll(io.LimitReader(context.Request.Body, maxLoggedBody))
			requestBody = string(raw)
			// คืน body กลับให้ handler ใช้ต่อ
			context.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
		}

		// ดัก response body
		writer := &bodyLogWriter{ResponseWriter: context.Writer}
		context.Writer = writer

		// ไปทำงานจริง
		context.Next()

		status := context.Writer.Status()
		route := context.FullPath()
		if route == "" {
			route = context.Request.URL.Path
		}
		module := moduleFromRoute(route)
		latency := time.Since(start)
		errMsg := context.Errors.ByType(gin.ErrorTypeAny).String()

		requestBody = sanitize(requestBody)
		responseBody := sanitize(writer.buffer.String())

		switch {
		case status >= 500:
			Errorf(module, "status=%d method=%s route=%s ip=%s latency=%s err=%s req=%s res=%s",
				status, context.Request.Method, route, context.ClientIP(), latency, sanitize(errMsg), requestBody, responseBody)
		case status >= 400:
			Warnf(module, "status=%d method=%s route=%s ip=%s latency=%s err=%s req=%s res=%s",
				status, context.Request.Method, route, context.ClientIP(), latency, sanitize(errMsg), requestBody, responseBody)
		default:
			Infof(module, "status=%d method=%s route=%s ip=%s latency=%s req=%s res=%s",
				status, context.Request.Method, route, context.ClientIP(), latency, requestBody, responseBody)
		}
	}
}
