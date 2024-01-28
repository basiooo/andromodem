package middleware

import (
	"context"
	adb "github.com/abccyz/goadb"
	"net/http"
)

func AdbClient(adbClient *adb.Adb) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			//lint:ignore SA1029 Using string as map key is allowed in this case.
			ctx := context.WithValue(request.Context(), "adbClient", adbClient)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
