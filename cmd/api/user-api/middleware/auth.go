package middleware

import (
	"context"
	"errors"
	"fmt"
	"ktrain/pkg/httputil"
	"ktrain/proto/pb"
	"net/http"
	"strings"
)

type ContextKey string
type dbTokenAuth struct {
	userClient pb.UserDMSServiceClient
}

func NewDBTokenAuth(userClient pb.UserDMSServiceClient) *dbTokenAuth {
	return &dbTokenAuth{
		userClient: userClient,
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
			var key ContextKey = "userID"
			ctx := context.WithValue(r.Context(), key, userID)
			fmt.Println("ok")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *dbTokenAuth) verifyToken(r *http.Request) (int64, error) {
	token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	fmt.Println(token)
	if token == "" {
		return 0, errors.New("empty token")
	}
	// result, err := m.userRepository.GetAuthToken(token)
	tokenReq := &pb.GetAuthTokenRequest{
		Token: token,
	}
	result, err := m.userClient.GetAuthToken(r.Context(), tokenReq)
	fmt.Println(result.AuthToken)
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
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *dbTokenAuth) verifyAdmin(r *http.Request) error {
	ctx := r.Context()
	// result, err := m.userRepository.GetUserByID(ctx.Value("userID").(int64))
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
