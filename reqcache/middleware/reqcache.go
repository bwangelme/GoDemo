package middleware

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"net/http"
)

const ContextKeyCache ContextKey = "reqCache"

func ReqCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, _ := lru.New[string, string](10000)
		ctx = CacheNewContext(ctx, c)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CacheNewContext(ctx context.Context, c *lru.Cache[string, string]) context.Context {
	ctx = context.WithValue(ctx, ContextKeyCache, c)
	return ctx
}

func CacheFromContext(ctx context.Context) (*lru.Cache[string, string], error) {
	reqCacheRaw := ctx.Value(ContextKeyCache)

	reqCache, ok := reqCacheRaw.(*lru.Cache[string, string])
	if !ok {
		return nil, fmt.Errorf("invalid req cache %v", reqCacheRaw)
	}

	return reqCache, nil
}
