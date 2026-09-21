package router

import (
	"database/sql"
	"net/http"

	"emergency-his/server/handler"
	"emergency-his/server/logger"
	"emergency-his/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string, db *sql.DB) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	gin.DefaultWriter = logger.Business.Writer()
	gin.DefaultErrorWriter = logger.Error.Writer()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(), middleware.RequestLog())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/api/login", handler.Login(db))
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	// 挂号模块路由先完成统一入口注册；具体业务按 Phase 2 分批接入。
	api.GET("/registration/searchPatients", handler.SearchPatients(db))
	api.GET("/registration/getPatientDetail", handler.GetPatientDetail(db))
	api.GET("/registration/getDepts", handler.GetDepts(db))
	api.GET("/registration/getDoctors", handler.GetDoctors(db))
	api.GET("/registration/getSchedules", handler.GetSchedules(db))
	api.POST("/registration/createRegistration", handler.CreateRegistration(db))
	api.GET("/registration/getEncounterList", handler.GetEncounterList(db))
	api.GET("/registration/getEncounterDetail", handler.GetEncounterDetail(db))
	api.POST("/registration/cancelEncounter", handler.CancelEncounter(db))
	return r
}
