package httputils

import (
	"context"

	"github.com/roaris/go-sns-api/models"
)

type contextKey string

const userContextKey contextKey = "user"

func SetUserToContext(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) models.User {
	user := ctx.Value(userContextKey)
	return user.(models.User)
}
