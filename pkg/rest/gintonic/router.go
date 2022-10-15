package gintonic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func AddHealthChecks(router *gin.Engine, database *gorm.DB) {
	checkRoutes := router.Group("/checks")
	{
		checkRoutes.GET("/healthz", healthz())
		checkRoutes.GET("/readiness", readiness(database))
	}
}

func healthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func readiness(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if db, err := database.DB(); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		} else {
			if err = db.Ping(); err != nil {
				c.JSON(http.StatusInternalServerError, nil)
				return
			}
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
