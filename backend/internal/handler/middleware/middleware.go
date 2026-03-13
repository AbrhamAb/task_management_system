package middleware

import (
	"context"
	"net/http"
	"strings"
	apperrors "task-management-backend/internal/constant/errors"
	"task-management-backend/internal/model/response"
	"task-management-backend/internal/module"
)

type ctxKey string

const userIDCtxKey ctxKey = "userID"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Auth(userModule *module.UserModule) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
				return
			}

			split := strings.SplitN(authHeader, " ", 2)
			if len(split) != 2 || !strings.EqualFold(split[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
				return
			}

			userID, err := userModule.ParseToken(split[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDCtxKey, userID)))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDCtxKey).(int64)
	return userID, ok
}
