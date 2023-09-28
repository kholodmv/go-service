// Package logger implements a simple logging package.
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// responseData struct include status and size response.
	responseData struct {
		status int // status response
		size   int // size response
	}

	// loggingResponseWriter struct include writer from http and response data.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Initialize function initialize logger.
func Initialize() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	return sugar
}

// Write function override standard function Write.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader function override standard function WriteHeader.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LoggerHandler that implements the logging.
func LoggerHandler(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic("cannot initialize zap")
		}
		defer logger.Sync()

		start := time.Now()
		duration := time.Since(start)

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		logger.Sugar().Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(f)
}
