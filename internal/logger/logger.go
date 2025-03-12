package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

}

type responseData struct {
	status     int
	size       int
	location   string
	respString string
}

type loggingResponseWriter struct {
	respData *responseData
	wr       http.ResponseWriter
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.wr.Write(b)
	r.respData.size += size // захватываем размер
	if size > 0 {
		r.respData.respString = string(b)
	}
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.wr.WriteHeader(statusCode)
	r.respData.status = statusCode // захватываем код статуса
	r.respData.location = r.wr.Header().Get("Location")
}

func (r *loggingResponseWriter) Header() http.Header {
	return r.wr.Header()
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := responseData{
			status: 0,
			size:   0,
		}

		rw := &loggingResponseWriter{
			respData: &rd,
			wr:       w,
		}

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(rw, r) // обслуживание оригинального запроса

		duration := time.Since(start)

		if method == http.MethodGet {
			sugar.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"status", rw.respData.status, // получаем перехваченный код статуса ответа
				"size", rw.respData.size,
				"Location", rw.respData.location,
			)
		} else //io.ReadAll(req.Body)
		if method == http.MethodPost {
			sugar.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"status", rw.respData.status, // получаем перехваченный код статуса ответа
				"size", rw.respData.size,
				"respString", rw.respData.respString,
			)
		} else {
			sugar.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"status", rw.respData.status, // получаем перехваченный код статуса ответа
				"size", rw.respData.size,
			)
		}
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
