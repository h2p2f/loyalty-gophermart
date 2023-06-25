package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log is httpserver logger
var Log *zap.Logger = zap.NewNop()

// responseData is struct for logging response data
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter is implementation of http.ResponseWriter
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write is implementation of http.ResponseWriter.Write
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader is implementation of http.ResponseWriter.WriteHeader
// It also sets status code to responseData
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)

}

// InitLogger is function for initializing logger
func InitLogger(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	//used for development logging
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	Log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

// WithLogging is middleware for logging
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{}
		loggedw := loggingResponseWriter{w, responseData}
		h.ServeHTTP(&loggedw, r)
		Log.Sugar().Infof("Request  - method: %s, url: %s, duration: %s", r.Method, r.URL.String(), time.Since(start))
		Log.Sugar().Infof("Request Info - Accept-Encoding: %s, Content-Encoding: %s", r.Header.Get("Accept-Encoding"), r.Header.Get("Content-Encoding"))
		Log.Sugar().Infof("Response - status: %d", responseData.status)
	}
	return http.HandlerFunc(logFn)
}
