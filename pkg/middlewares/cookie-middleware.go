package middleware

import (
	"context"
	"net/http"
)

type cookieContextKey struct{}

func WithSetCookie(ctx context.Context, cookie *http.Cookie) context.Context {
	return context.WithValue(ctx, cookieContextKey{}, cookie)
}

func SetCookieFromContext(ctx context.Context, w http.ResponseWriter) {
	cookie, ok := ctx.Value(cookieContextKey{}).(*http.Cookie)
	if ok {
		http.SetCookie(w, cookie)
	}
}
