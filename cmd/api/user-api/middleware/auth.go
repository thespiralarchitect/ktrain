package middleware

import (
	"context"
	"errors"
	"ktrain/pkg/httputil"
	"ktrain/proto/pb"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ContextKey string
type dbTokenAuth struct {
	userClient pb.UserDMSServiceClient
	logger     *zap.SugaredLogger
}

func NewDBTokenAuth(userClient pb.UserDMSServiceClient, logger *zap.SugaredLogger) *dbTokenAuth {
	return &dbTokenAuth{
		userClient: userClient,
	}
}

func (m *dbTokenAuth) Handle() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := m.verifyToken(r)
			if err != nil {
				m.logger.Errorw("Error verify token", "error", err)
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}
			var key ContextKey = "userID"
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
	tokenReq := &pb.GetAuthTokenRequest{
		Token: token,
	}
	result, err := m.userClient.GetAuthToken(r.Context(), tokenReq)
	if err != nil {
		return 0, errors.New("invalid token")
	}
	return result.AuthToken.UserId, nil
}

func (m *dbTokenAuth) HandleAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := m.verifyAdmin(r)
			if err != nil {
				m.logger.Errorw("Error verify admin", "error", err)
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *dbTokenAuth) verifyAdmin(r *http.Request) error {
	ctx := r.Context()
	getUserReq := &pb.GetUserByIDRequest{
		Id: ctx.Value(ContextKey("userID")).(int64),
	}
	result, err := m.userClient.GetUserByID(r.Context(), getUserReq)
	if err != nil {
		return err
	}
	if !result.User.IsAdmin {
		return errors.New("Permission denied")
	}
	return nil
}
