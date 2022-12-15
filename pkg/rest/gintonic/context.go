package gintonic

import (
	"context"
	"github.com/google/uuid"
)

type UserRole int8

const (
	RoleBasicUser UserRole = iota
	RoleAdmin

	CtxKeyUserID = "userID"
	CtxKeyRole   = "userRole"
	CtxKeyApp    = "appID"
)

func GetAppID(ctx context.Context) string {
	return ctx.Value(CtxKeyApp).(string)
}

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(CtxKeyUserID).(uuid.UUID)
}

func GetUserRole(ctx context.Context) UserRole {
	return ctx.Value(CtxKeyRole).(UserRole)
}

func IsBasicUser(ctx context.Context) bool {
	return GetUserRole(ctx) == RoleBasicUser
}

func IsAdmin(ctx context.Context) bool {
	return GetUserRole(ctx) == RoleAdmin
}
