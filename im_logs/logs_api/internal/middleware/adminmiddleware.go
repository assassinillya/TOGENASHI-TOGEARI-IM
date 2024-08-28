package middleware

import (
	"fmt"
	"im_server/common/response"
	"k8s.io/kube-openapi/pkg/validation/errors"
	"net/http"
)

type AdminMiddleware struct {
}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只有角色为1的用户才能调用
		role := r.Header.Get("Role")
		userid := r.Header.Get("User-ID")
		fmt.Println(userid)
		if role != "1" {
			response.Response(r, w, nil, errors.New(200, "角色权限不足"))
			return
		}
		next(w, r)
	}
}
