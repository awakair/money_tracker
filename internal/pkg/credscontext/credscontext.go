package credscontext

import (
	"context"
	"money_tracker/internal/domain"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	userRoleKey contextKey = "user_role"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithUserRole(ctx context.Context, role domain.UserRole) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func UserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

func UserRoleFromContext(ctx context.Context) domain.UserRole {
	role, _ := ctx.Value(userRoleKey).(domain.UserRole)
	return role
}
