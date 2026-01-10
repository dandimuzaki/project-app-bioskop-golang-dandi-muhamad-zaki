package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/project-app-bioskop-golang/internal/data/entity"
	"github.com/project-app-bioskop-golang/pkg/utils"
	"go.uber.org/zap"
)

func (mw *MiddlewareCustom) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			token := strings.TrimSpace(strings.Replace(auth, "Bearer", "", 1))
			userID, err := mw.Usecase.AuthUsecase.ValidateToken(token)
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "invalid token", err)
				return
			}

			user, err := mw.Usecase.UserUsecase.GetByID(*userID)
			if err != nil {
				utils.ResponseFailed(w, http.StatusUnauthorized, "user not found", err)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (middlewareCostume *MiddlewareCustom) RequirePermission(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(entity.User)
			if !ok {
				middlewareCostume.Log.Error("Error retrieve user info", zap.Error(errors.New("error retrieve user")))
				utils.ResponseFailed(w, http.StatusUnauthorized, "invalid user", errors.New("error retrieve user"))
				return
			}

			isAllowed := false
			for _, role := range roles {
				if user.Role == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				utils.ResponseFailed(w, http.StatusForbidden, "invalid role", errors.New("invalid role"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
