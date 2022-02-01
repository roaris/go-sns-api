package httputils

import "context"

type contextKey string

const userIDContextKey contextKey = "userID"

func SetUserIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetUserIDFromContext(ctx context.Context) int64 {
	userID := ctx.Value(userIDContextKey)
	return userID.(int64)
}
