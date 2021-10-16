package middleware

import (
	"context"
	"errors"
	"ktrain/pkg/httputil"
	"ktrain/pkg/tokens"
	"net/http"
	"strings"
)

func (m *dbTokenAuth) verifyJWTToken(r *http.Request) (int64, error) {
	jwtToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if jwtToken == "" {
		return 0, errors.New("empty token")
	}
	tokenAuth, err := tokens.ParseJWT(jwtToken)
	if err != nil {
		return 0, err
	}
	claims, ok := tokenAuth.Claims.(*tokens.MyClaims)
	if !ok {
		return 0, errors.New("error convert claims")
	}
	return claims.UserID, nil
}

func (m *dbTokenAuth) HandleJWT() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := m.verifyJWTToken(r)
			if err != nil {
				httputil.RespondError(w, http.StatusForbidden, err.Error())
				return
			}
			var key ContextKey = "userID"
			ctx := context.WithValue(r.Context(), key, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
