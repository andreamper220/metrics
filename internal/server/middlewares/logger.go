package middlewares

import (
	"net/http"
	"time"

	"github.com/andreamper220/metrics.git/internal/logger"
)

type (
	responseData struct {
		size int
		code int
	}
	loggingResponseWriter struct {
		w    http.ResponseWriter
		data *responseData
	}
)

func (lw *loggingResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lw.w.Write(data)
	if err != nil {
		return 0, err
	}

	lw.data.size += size
	return size, nil
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.w.WriteHeader(code)

	lw.data.code = code
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lw := loggingResponseWriter{
			w: w,
			data: &responseData{
				size: 0,
				code: 0,
			},
		}

		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h(&lw, r)

		duration := time.Since(start).Milliseconds()

		logger.Log.Infof("REQUEST  | URI: %s, Method: %s, Duration: %dms", uri, method, duration)
		logger.Log.Infof("RESPONSE | Status: %d, Size: %d", lw.data.code, lw.data.size)
	}
}
