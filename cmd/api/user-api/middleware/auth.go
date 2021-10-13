package middleware

import (
	"context"
	"errors"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
	"strings"
)

type keyUserID string
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
			var key keyUserID = "userID"
			ctx := context.WithValue(r.Context(), key, userID)
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

func (m *dbTokenAuth) HandleAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := m.verifyAdmin(r)
			if err != nil {
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *dbTokenAuth) verifyAdmin(r *http.Request) error {
	ctx := r.Context()
	result, err := m.userRepository.GetUserByID(ctx.Value("userID").(int64))
	if err != nil {
		return err
	}
	if !result.IsAdmin {
		return errors.New("Permission denied")
	}
	return nil
}
