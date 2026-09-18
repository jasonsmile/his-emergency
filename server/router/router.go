package router

import (
	"database/sql"
	"net/http"

	"emergency-his/server/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string, db *sql.DB) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/api/v1/patients", handler.ListPatients(db))
	return r
}
