package gintonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"strconv"
)

const (
	HeaderUserID = "User-Id"
	HeaderRole   = "Role"
	HeaderApp    = "App-Id"
)

func UserScoped(c *gin.Context, ctx context.Context) (context.Context, bool) {
	log := ctxlogrus.Extract(ctx)
	// Extract user role and select enum
	userRoleStr := c.GetHeader(HeaderRole)
	if userRoleStr == "" {
		log.Warnf("request is missing header '%v'", HeaderRole)
		userRoleStr = "0" // Aka Basic user
	}
	userRole, err := strconv.ParseInt(userRoleStr, 10, 8)
	if err != nil {
		log.Errorf("'%v' is not int8", HeaderRole)
		return ctx, false
	}
	userRoleEnum := UserRole(userRole)
	switch userRoleEnum {
	case RoleBasicUser:
		fallthrough
	case RoleAdmin:
		break
	default:
		log.Errorf("request is missing header '%v'", HeaderRole)
		return ctx, false
	}

	// Extract user id and parse as uuid
	userID := c.GetHeader(HeaderUserID)
	if userID == "" {
		log.Errorf("request is missing header '%v'", HeaderUserID)
		return ctx, false
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Debugf("header '%v' is not in uuid format %v", HeaderUserID, userID)
		return ctx, false
	}

	// Attach data to context
	ctx = context.WithValue(ctx, CtxKeyUserID, userUUID)
	ctx = context.WithValue(ctx, CtxKeyRole, userRoleEnum)
	return ctx, true
}

func ApplicationScoped(c *gin.Context, ctx context.Context) (context.Context, bool) {
	log := ctxlogrus.Extract(ctx)

	// Extract user id and parse as uuid
	appID := c.GetHeader(HeaderApp)
	if appID == "" {
		log.Errorf("request is missing header '%v'", HeaderApp)
		return ctx, false
	}

	// Attach data to context
	ctx = context.WithValue(ctx, CtxKeyApp, appID)
	return ctx, true
}
