package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/api/middleware"
	v1 "monitorlogs/internal/api/v1"
	"monitorlogs/internal/api/v2"
	"monitorlogs/internal/config"
	"monitorlogs/internal/db"
	"monitorlogs/internal/logreader"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"strings"
)

func init() {
	err := config.SetEnv()
	if err != nil {
		tools.LogErr(erx.New(err))
		return
	}
	tools.LogInfo("Env is successfully set")
}

// IMPORTANT! Build with --tags "fts5" For Full Text Search Support

func main() {
	function := flag.String("f", "default", "Specify one of the commands: initUsersDb -p, initLogsDb -p, read -fn, server")
	path := flag.String("p", "./data/", "Specify PATH (to Database Folder or to Folder with logs)")
	fileName := flag.String("fn", `.\data\debug.log`, "Specify Filename of file to read (with single backslash)")
	flag.Parse()

	switch *function {
	case "initUsersDb":
		err := db.InitUsersDb(*path)
		if err == nil {
			tools.LogInfo("Database Users initialised successfully")
		}
		if err != nil {
			if strings.Contains(err.Error(), "605") {
				tools.LogWarn("Unable to init Logs Database")
				err = nil
			} else {
				tools.LogErr(erx.New(err))
				return
			}
		}
	case "initLogsDb":
		err := db.InitLogsDb(*path)
		if err == nil {
			tools.LogInfo("Database Logs initialised successfully")
		}
		if err != nil {
			if strings.Contains(err.Error(), "605") {
				tools.LogWarn("Unable to init Logs Database: Database already exists")
				err = nil
			} else {
				tools.LogErr(erx.New(err))
				return
			}
		}
	case "read":
		err := logreader.Read(*fileName)
		if err != nil {
			tools.LogErr(erx.New(err))
		}
	case "readFolder":
		logreader.ReadFolder(*path)
	case "server":
		r := gin.Default()

		conf, err := config.GetConfig()
		if err != nil {
			erx.New(err)
			return
		}

		logfiles, err := db.GetLogsFilesInfo()
		if err == sql.ErrNoRows {
			tools.LogWarn("NO LOGS INFO FOUND!")
			logfiles = nil
		}
		if err != nil {
			tools.LogErr(erx.New(err))
			erx.New(err)
			return
		}
		err = config.SetLogfilesEnv(logfiles)
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}

		go logreader.ReadCycle(conf.Logs.ReadCycle, conf.Logs.Path)
		r.Use(static.Serve("/", static.LocalFile(conf.Templates.Path, true)))

		api1 := r.Group("/v1")
		{
			//api1.GET("/", v1.ShowMain)
			api1.GET("/auth", v1.Auth)
			api1.GET("/loginAttempt", v1.ShowLoginPage)
			api1.GET("/logoutAttempt", v1.Logout)
			api1.POST("/register1", v1.Register)
			api1.GET("/registration", v1.ShowRegistrationPage)
			api1.POST("/loginAttempt", v1.Login)
			api1.POST("/unblock", v1.Unblock)
		}

		api2 := r.Group("/v2")
		{
			api2.GET("/authAttempt", v2.Auth)
			api2.POST("/loginAttempt", v2.Login)
			api2.GET("/logoutAttempt", v2.Logout)
			api2.POST("/registrationAttempt", v2.Register)
		}

		authorized := r.Group("/v2/private")
		authorized.Use(middleware.AuthJWT())
		{
			authorized.POST("/block", v2.Block)
			authorized.POST("/unblock", v2.Unblock)
			authorized.POST("/getLogsBySession", v2.GetLogsBySession)
			authorized.POST("/getLogsByDate", v2.GetLogsByDate)
			authorized.POST("/getLogsBySessionWithLimit", v2.GetLogsBySessionWithLimit)
			authorized.POST("/getLogsByDateWithLimit", v2.GetLogsByDateWithLimit)
			authorized.POST("/getErrorsBySessionWithLimit", v2.GetErrorsBySessionWithLimit)
			authorized.POST("/getErrorsByDateWithLimit", v2.GetErrorsByDateWithLimit)
			authorized.POST("/findLogs", v2.FindLogs)
			authorized.POST("/findLogsWithLimit", v2.FindLogsWithLimit)
			authorized.POST("/getLogById", v2.GetLogById)
			authorized.POST("/getLogsSessions", v2.GetLogsSessions)
			authorized.POST("/getLogsServiceInfo", v2.GetLogsServiceInfo)
			authorized.GET("/getLogsFilenames", v2.GetLogsFilenames)
		}

		//r.POST("/hsmauth/v1/block", v1.Block)
		//r.POST("/hsmauth/v1/unblock", v1.Unblock)
		//r.POST("/hsmauth/v1/exitAll", v1.ExitAll)

		port := fmt.Sprintf(":%v", conf.Server.Port)
		if conf.TLS.Enable {

			SSLCRT := conf.TLS.Certificate
			SSLKEY := conf.TLS.Key
			tools.LogInfo("Starting TLS server")
			err = r.RunTLS(port, SSLCRT, SSLKEY)
			if err != nil {
				tools.LogErr(erx.New(err))
				return
			}
		}
		// TODO Make Dev/Prod separator
		tools.LogWarn("Starting WITHOUT TLS server")
		corsConf := cors.DefaultConfig()
		corsConf.AllowOrigins = []string{"http://localhost:3000"}
		corsConf.AllowCredentials = true
		corsConf.AllowHeaders = []string{"Fingerprint", "X-Requested-With", "content-type", "Authorization", "Set-Cookie"}
		corsConf.AllowMethods = []string{"GET", "POST"}
		r.Use(cors.New(corsConf))
		err = r.Run(port)
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}

	default:
		fmt.Println("Expected flag (-f)")
		flag.PrintDefaults()
	}
}
