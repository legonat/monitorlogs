package db

import (
	"database/sql"
	"monitorlogs/internal/config"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATE_TABLE_USERS = `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
	salt VARCHAR(255) NOT NULL, 
	create_at INTEGER NULL, 
	blocked BOOL NULL,
	try_count INTEGER NULL,
	blocked_at INTEGER NULL,
	deleted BOOL NULL);`

	CREATE_TABLE_STATISTIC = `
	CREATE TABLE statistics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event INTEGER NOT NULL,
	create_at INTEGER NOT NULL, 
	ip VARCHAR(15) NOT NULL,
	details VARCHAR(255) NULL);
	`

	CREATE_TABLE_EVENTS = `
	CREATE TABLE events (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(100) NOT NULL, 
	caption VARCHAR(255) NOT NULL);
	`

	CREATE_TABLE_SESSIONS = `
	CREATE TABLE sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    login VARCHAR(100) NOT NULL,
    refreshToken INTEGER NOT NULL,
    ua VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(255) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    expiresIn INTEGER NOT NULL,
    createdAt INTEGER NOT NULL);
	`

	CREATE_TABLE_LOG_FILES = `
	CREATE TABLE logfiles (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	logfileName TEXT NOT NULL,
	fileLength INTEGER NOT NULL,
	lastSessionDate INTEGER NULL,
	previousDate INTEGER NULL,
	sessionCount INTEGER NOT NULL);
	`

)

func NewSqliteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitUsersDb(path string) error {

	pathDb := path + "/auth.db"
	cfg, err := config.GetConfig()
	if err != nil {
		tools.LogErr(err)
		return err
	}
	(*cfg).UsersDB.PathDb = pathDb

	f, err := os.Stat(pathDb)
	if err != nil && !os.IsNotExist(err) {
		return erx.New(err)
	}
	if f != nil {
		return erx.NewError(605, "Data Base is already exist")
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return erx.New(err)
		}

		_, err = os.Create(pathDb)
		if err != nil {
			return erx.New(err)
		}

		db, err := sql.Open("sqlite3", pathDb)
		if err != nil {
			return erx.New(err)
		}

		defer db.Close()

		_, err = db.Exec(CREATE_TABLE_USERS)

		if err != nil {
			return erx.New(err)
		}
		_, err = db.Exec(CREATE_TABLE_STATISTIC)
		if err != nil {
			return erx.New(err)
		}

		_, err = db.Exec(CREATE_TABLE_EVENTS)
		if err != nil {
			return erx.New(err)
		}
		for _, v := range models.Events {
			_, err = db.Exec(INSERT_EVENT, v.Name, v.Caption)
			if err != nil {
				return erx.New(err)
			}
		}
		_, err = db.Exec(CREATE_TABLE_SESSIONS)
		if err != nil {
			return erx.New(err)
		}

		err = cfg.RewriteConfig()
		if err != nil {
			tools.LogErr(erx.New(err))
			return erx.New(err)
		}
	}

	return err
}

func InitLogsDb(path string) error {

	pathDb := path + "/logs.db"
	cfg, err := config.GetConfig()
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	(*cfg).Logs.PathDb = pathDb

	f, err := os.Stat(pathDb)
	if err != nil && !os.IsNotExist(err) {
		return erx.New(err)
	}
	if f != nil && f.Size() > 0 {
		return erx.NewError(605, "Data Base already exists")
	}
	if os.IsNotExist(err) || f.Size() == 0 {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return erx.New(err)
		}

		_, err = os.Create(pathDb)
		if err != nil {
			return erx.New(err)
		}

		db, err := sql.Open("sqlite3", pathDb)
		if err != nil {
			return erx.New(err)
		}

		defer db.Close()

		_, err = db.Exec(CREATE_TABLE_LOG_FILES)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}

		//cfg.Logs.SessionCount = 1
		//cfg.Logs.Length = 0
		//cfg.Logs.LastSessionDate = ""
		//cfg.Logs.PreviousDate = ""
		err = cfg.RewriteConfig()
		if err != nil {
			tools.LogErr(erx.New(err))
			return erx.New(err)
		}
	}
	return nil
}
