package middleware

import (
	"context"
	"encoding/json"
	"go-auth-with-chi/domain"
	"go-auth-with-chi/dto"
	"go-auth-with-chi/ioc"
	"go-auth-with-chi/utils"
	"net/http"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// AuthMiddleware :
// case 1; invalid token / expired token
// case 2; no token given
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			header := r.Header.Get("Authorization")

			// case 2; no token given
			if header == "" {
				ctx := context.WithValue(r.Context(), UserCtxKey, utils.ErrNoTokenGiven)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			// case 1; invalid token
			tokenStr := header
			id, err := utils.ParseToken(tokenStr)
			if err != nil {
				var data dto.Msg
				data.Error = utils.ErrInvalidToken.Error()
				if err == utils.ErrExpiredToken {
					data.Error = err.Error()
					data.Result = "fail"
					jsonData, _ := json.Marshal(data)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(jsonData)
					return
				}
				data.Result = "fail"
				jsonData, _ := json.Marshal(data)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(jsonData)
				return
			}

			// case 1; claim not matched with user data
			user, err := ioc.Repo.Users.Get(id)
			if err != nil {
				var data dto.Msg
				data.Error = utils.ErrInvalidToken.Error()
				data.Result = "fail"
				jsonData, _ := json.Marshal(data)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(jsonData)
				return
			}

			// add result to context
			ctx := context.WithValue(r.Context(), UserCtxKey, user)

			// call next with new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user or errors from the context. REQUIRES AuthMiddleware to have run.
func ForContext(ctx context.Context) (*domain.UserDAO, error) {
	raw, ok := ctx.Value(UserCtxKey).(*domain.UserDAO)
	if ok {
		return raw, nil
	}

	errMsg, _ := ctx.Value(UserCtxKey).(error)
	return nil, errMsg
}