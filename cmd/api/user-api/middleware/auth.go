package middleware

import (
	"context"
	"errors"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
	"strings"
)

type dbTokenAuth struct {
	userRepository repository.IUserRepository
}

func NewDBTokenAuth(userRepository repository.IUserRepository) *dbTokenAuth {
	return &dbTokenAuth{
		userRepository: userRepository,
	}
}

func (m *dbTokenAuth) Handle() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := m.verifyToken(r)
			if err != nil {
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func (m *dbTokenAuth) verifyToken(r *http.Request) (int64, error) {
	token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if token == "" {
		return 0, errors.New("empty token")
	}

	result, err := m.userRepository.GetAuthToken(token)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return result.UserID, nil
}
