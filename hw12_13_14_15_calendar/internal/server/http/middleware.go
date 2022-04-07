package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func addLoggingMiddleware(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w}
		clientIP := r.RemoteAddr
		url := r.URL.Path
		httpProtocol := r.Proto
		method := r.Method
		userAgent := strings.Split(r.UserAgent(), " ")[0]

		start := time.Now()
		next.ServeHTTP(&sw, r)
		duration := time.Since(start)

		logger.Info(
			fmt.Sprintf("http request has been made...\nClientIP:%s;\nMethod:%s;\nURL:%s;\nHttpProtocol:%s;\nStatusCode:%d;"+
				"\nContentLength:%d;\nLatency:%s;\nUser_agent:%s",
				clientIP, method, url, httpProtocol, sw.status, sw.length, duration, userAgent),
		)
	})
}
