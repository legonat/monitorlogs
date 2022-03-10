package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"monitorlogs/internal/config"
	"monitorlogs/internal/db"
	"monitorlogs/internal/handler"
	"monitorlogs/internal/server"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
		usersDb, err := db.NewSqliteDB(os.Getenv("USERS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		defer usersDb.Close()
		logfilesDb, err := db.NewSqliteDB(os.Getenv("LOGS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		defer logfilesDb.Close()

		repos := db.NewRepository(logfilesDb, usersDb)
		err = repos.Read(*fileName)
		if err != nil {
			tools.LogErr(erx.New(err))
		}
	case "readFolder":
		usersDb, err := db.NewSqliteDB(os.Getenv("USERS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		defer usersDb.Close()
		logfilesDb, err := db.NewSqliteDB(os.Getenv("LOGS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		defer logfilesDb.Close()

		repos := db.NewRepository(logfilesDb, usersDb)
		repos.ReadFolder(*path)
	case "server":

		conf, err := config.GetConfig()
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}

		usersDb, err := db.NewSqliteDB(os.Getenv("USERS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		logfilesDb, err := db.NewSqliteDB(os.Getenv("LOGS_PATH_DB"))
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}

		repos := db.NewRepository(logfilesDb, usersDb)
		mainHandler := handler.NewHandler(repos)

		logfiles, err := repos.GetLogsFilesInfo()
		if err == sql.ErrNoRows {
			tools.LogWarn("NO LOGS INFO FOUND!")
			logfiles = nil
		}
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}
		err = config.SetLogfilesEnv(logfiles)
		if err != nil {
			tools.LogErr(erx.New(err))
			return
		}

		go repos.ReadCycle(conf.Logs.ReadCycle, conf.Logs.Path)

		port := fmt.Sprintf("%v", conf.Server.Port)

		srv := new(server.Server)

		if conf.TLS.Enable {
			SSLCRT := conf.TLS.Certificate
			SSLKEY := conf.TLS.Key
			tools.LogInfo("Starting TLS server")
			go func() {
				if err := srv.RunTLS(port, mainHandler.InitRoutes(), SSLCRT, SSLKEY); err != nil {
					tools.LogErr(erx.New(err))
				}
			}()
		} else {
			go func() {
				if err := srv.Run(port, mainHandler.InitRoutes()); err != nil {
					tools.LogErr(erx.New(err))
				}
			}()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		tools.LogInfo("Monitorlogs Shutting Down")

		if err := srv.Shutdown(context.Background()); err != nil {
			tools.LogErr(erx.New(err))
		}

		if err := usersDb.Close(); err != nil {
			tools.LogErr(erx.New(err))
		}

		if err := logfilesDb.Close(); err != nil {
			tools.LogErr(erx.New(err))
		}



	default:
		fmt.Println("Expected flag (-f)")
		flag.PrintDefaults()
	}
}
