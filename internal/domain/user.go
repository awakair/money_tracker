package domain

import "time"

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CanViewResource(userID string, role UserRole, resourceOwnerID string) bool {
	return userID == resourceOwnerID || role == UserRoleAdmin
}

func CanEditResource(userID string, role UserRole, resourceOwnerID string) bool {
	return userID == resourceOwnerID
}
