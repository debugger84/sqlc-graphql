package auth

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

var ErrUnauthenticated = errors.New("user is not authenticated: missing Authorization header")

type contextKey string

func setUserId(ctx context.Context, currentUserId *uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKey("CurrentUserId"), currentUserId)
}

func GetCurrentUserId(ctx context.Context) *uuid.UUID {
	if value := ctx.Value(contextKey("CurrentUserId")); value != nil {
		currentUserId, ok := value.(*uuid.UUID)
		if !ok {
			return nil
		}
		return currentUserId
	}
	return nil
}

// AuthenticateMiddleware It is a stub for the authentication mechanics.
// It is used only to show an example of post and comment creation as close as possible to the real usage.
// It rewrites the request with the current user id
func AuthenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader == "" {
			next(writer, request)
			return
		}
		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 {
			next(writer, request)
			return
		}
		if parts[0] != "Bearer" {
			next(writer, request)
			return
		}

		ctx := request.Context()
		userId, err := uuid.FromString(parts[1])
		if err != nil {
			next(writer, request)
			return
		}
		ctx = setUserId(ctx, &userId)
		request = request.WithContext(ctx)
		next(writer, request)
	}
}
