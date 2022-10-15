package gintonic

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"net/http"
	"strconv"
)

const (
	HeaderUserID = "UniFyi-User-Id"
	HeaderRole   = "UniFyi-Role"
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
		c.JSON(http.StatusInternalServerError, nil)
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
