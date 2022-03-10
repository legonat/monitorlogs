package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/api/middleware"
	"monitorlogs/internal/config"
	"monitorlogs/internal/db"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
)

type Handler struct {
	repository *db.Repository
}

func NewHandler(repository *db.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	conf, err := config.GetConfig()
	if err != nil {
		tools.LogErr(erx.New(err))
	}

	r.Use(static.Serve("/", static.LocalFile(conf.Templates.Path, true)))


	api2 := r.Group("/v2")
	{
		api2.GET("/authAttempt", h.Auth)
		api2.POST("/loginAttempt", h.Login)
		api2.GET("/logoutAttempt", h.Logout)
		api2.POST("/registrationAttempt", h.Register)
	}

	authorized := r.Group("/v2/private")
	authorized.Use(middleware.AuthJWT())
	{
		authorized.POST("/block", h.Block)
		authorized.POST("/unblock", h.Unblock)
		authorized.POST("/getLogsBySession", h.GetLogsBySession)
		authorized.POST("/getLogsByDate", h.GetLogsByDate)
		authorized.POST("/getLogsBySessionWithLimit", h.GetLogsBySessionWithLimit)
		authorized.POST("/getLogsByDateWithLimit", h.GetLogsByDateWithLimit)
		authorized.POST("/getErrorsBySessionWithLimit", h.GetErrorsBySessionWithLimit)
		authorized.POST("/getErrorsByDateWithLimit", h.GetErrorsByDateWithLimit)
		authorized.POST("/findLogs", h.FindLogs)
		authorized.POST("/findLogsWithLimit", h.FindLogsWithLimit)
		authorized.POST("/getLogById", h.GetLogById)
		authorized.POST("/getLogsSessions", h.GetLogsSessions)
		authorized.POST("/getLogsServiceInfo", h.GetLogsServiceInfo)
		authorized.GET("/getLogsFilenames", h.GetLogsFilenames)
	}

	// TODO Make Dev/Prod separator
	tools.LogWarn("Starting WITHOUT TLS server")
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"http://localhost:3000"}
	corsConf.AllowCredentials = true
	corsConf.AllowHeaders = []string{"Fingerprint", "X-Requested-With", "content-type", "Authorization", "Set-Cookie"}
	corsConf.AllowMethods = []string{"GET", "POST"}
	r.Use(cors.New(corsConf))

	return r
}