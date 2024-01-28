package middleware

import (
	"context"
	adb "github.com/abccyz/goadb"
	"net/http"
)

type contextKey string

const AdbClientKey contextKey = "adbClient"

func AdbClient(adbClient *adb.Adb) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), AdbClientKey, adbClient)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
