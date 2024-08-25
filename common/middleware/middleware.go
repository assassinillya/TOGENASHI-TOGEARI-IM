package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
	Body []byte
}

func (w *Writer) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func LogActionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := httpx.GetRemoteAddr(r)

		ctx := context.WithValue(r.Context(), "clientIP", clientIP)
		ctx = context.WithValue(ctx, "userID", r.Header.Get("User-ID"))
		next(w, r.WithContext(ctx))

	}

}
