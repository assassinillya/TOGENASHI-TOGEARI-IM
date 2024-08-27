package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"im_server/common/log_stash"
	"net/http"
)

func LogActionMiddleware(pusher *log_stash.Pusher) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			clientIP := httpx.GetRemoteAddr(r)

			// 设置入参
			pusher.SetRequest(r)
			pusher.SetHeaders(r)

			ctx := context.WithValue(r.Context(), "clientIP", clientIP)
			ctx = context.WithValue(ctx, "userID", r.Header.Get("User-ID"))
			next(w, r.WithContext(ctx))
			// 设置响应
			pusher.SetResponse(w)
		}
	}
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := httpx.GetRemoteAddr(r)
		ctx := context.WithValue(r.Context(), "clientIP", clientIP)
		ctx = context.WithValue(ctx, "userID", r.Header.Get("User-ID"))
		next(w, r.WithContext(ctx))

	}

}
