package httputils

import "context"

type contextKey string

const userIDContextKey contextKey = "userID"

func SetUserIDToContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetUserIDFromContext(ctx context.Context) int {
	userID := ctx.Value(userIDContextKey)
	return userID.(int)
}
