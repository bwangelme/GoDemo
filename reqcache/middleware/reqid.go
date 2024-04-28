package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"context"
	"github.com/google/uuid"
)

var (
	ReqidHeader = "X-Request-Id"
	logger      = logrus.New()
)

type ContextKey string

const ContextKeyRequestID ContextKey = "requestID"

// FromContext 获取存储在 ctx 中的 req id
func FromContext(ctx context.Context) (string, error) {
	reqIDRaw := ctx.Value(ContextKeyRequestID)

	reqID, ok := reqIDRaw.(string)
	if !ok {
		return "", fmt.Errorf("invalid reqid %v", reqIDRaw)
	}

	return reqID, nil
}

// NewContext
// 将 reqid 设置到 context 中
func NewContext(ctx context.Context, id string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyRequestID, id)
	return ctx
}

// ReqIDMiddleware
// 在请求上下文中设置 request id，如果设置失败，会终止请求
func ReqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := getReqIDFromReq(r)
		if err != nil {
			logger.Warning("generate req id error %v", err)
			http.Error(w, fmt.Sprintf("get reqid error %v", err), http.StatusInternalServerError)
			return
		}
		r = setRequestID(r, id)
		next.ServeHTTP(w, r)
	})
}

// getReqIDFromReq
// 从 request 中获取 req id，如果获取失败，则使用 uuid 作为 request id
// 上述两个方法都失败后，会返回 err，终止当前请求
func getReqIDFromReq(r *http.Request) (string, error) {
	id := r.Header.Get(ReqidHeader)
	if id != "" {
		return id, nil
	}
	_uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return _uuid.String(), nil

}

func setRequestID(req *http.Request, id string) *http.Request {
	ctx := req.Context()
	ctx = NewContext(ctx, id)
	return req.WithContext(ctx)
}
