package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"monitorlogs/internal/config"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//func init ()  {
//	tools.LogrusWithParams("info", nil, "init in db")
//	conf := config.New()
//	path = conf.Auth.PathDb
//}

func getUsersDbPath() string {
	//conf := config.New()
	//return conf.Auth.PathDb
	return os.Getenv("USERS_PATH_DB")
}
func getLogsDbPath() string {
	//conf := config.New()
	//return conf.Auth.PathDb
	return os.Getenv("LOGS_PATH_DB")
}

func GetTime() int64 {
	return time.Now().Unix()
}

func GetExpTime(days int) int64 {
	return time.Now().Add(time.Hour * 24 * time.Duration(days)).Unix()
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)

	return salt, erx.New(err)
}

func PasswordHash(password []byte, salt []byte) []byte {
	a := append(salt, password...)
	h := sha256.New()
	h.Write(a)
	return h.Sum(nil)
}

var logsNotFoundStruct models.LogStruct = models.LogStruct{Id: 1, ServiceInfo: "Not Found", Description: "Logs not found"}

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

	CREATE_TABLE_LOGS = `
	CREATE TABLE %v (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sessionId INTEGER NOT NULL, 
	date INTEGER NOT NULL,
	serviceInfo VARCHAR(255) NULL,
	description TEXT NULL);
	`

	CREATE_TABLE_LOGS_FTS5 = `
	CREATE VIRTUAL TABLE %v USING FTS5(
	id,
	sessionId, 
	date,
	serviceInfo,
	description);
	`

	CREATE_TABLE_ERRORS = `
	CREATE TABLE %v (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	logId INTEGER NOT NULL,
	sessionId INTEGER NOT NULL, 
	date INTEGER NOT NULL,
	serviceInfo VARCHAR(255) NULL,
	description TEXT NULL);
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

	CREATE_TABLE_SESSION_DATES = `
	CREATE TABLE %v (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	dates INTEGER NOT NULL);
	`

	WRITE_REFRESH_TOKEN = `INSERT INTO sessions (login, refreshToken, ua, fingerprint, ip, expiresIn, createdAt) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	GET_REFRESH_SESSION = `SELECT login, refreshToken, fingerprint, expiresIn, createdAt FROM sessions WHERE refreshToken = $1;`

	GET_ALL_REFRESH_SESSIONS = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE login Like $1 ORDER BY createdAt ASC;`

	GET_ALL_REFRESH_SESSIONS_SORTED = `SELECT refreshToken, fingerprint FROM sessions WHERE login Like $1 ORDER BY createdAt ASC;`

	FIND_REFRESH_TOKEN = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE fingerprint = $1;`

	FIND_REFRESH_SESSION = `SELECT login, refreshToken, fingerprint, ip FROM sessions WHERE refreshToken = $1;`

	DELETE_TOKEN = `DELETE FROM sessions WHERE refreshToken = $1;`

	//INVALIDATE_TOKEN = `UPDATE sessions SET valid = 0 WHERE refreshToken = $1;`

	INSERT_EVENT = `INSERT INTO events (name, caption)VALUES ($1, $2);`

	WRITE_STAT = `INSERT INTO statistics (event, create_at, ip, details) VALUES ($1, $2 ,$3, $4);`

	GET_USER = `SELECT login FROM users WHERE login LIKE $1;`

	WRITE_USER = `INSERT INTO users (login, password, salt, create_at, blocked, try_count, blocked_at, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	GET_USER_PASSWORD = `SELECT password, salt, blocked, try_count FROM users WHERE login LIKE ?;`

	WRITE_TRY_COUNT = `UPDATE users SET try_count = $1 WHERE login LIKE $2 `

	BLOCK_USER = `UPDATE users SET blocked = 1, blocked_at = $1 WHERE login LIKE $2`

	GET_BLOCKED = `SELECT blocked FROM users WHERE login LIKE $1;`

	UNBLOCK_USER = `UPDATE users SET blocked = 0, try_count = 5, blocked_at = 0 WHERE login LIKE $1`

	GET_LOG_BY_ID = `SELECT * FROM %v WHERE id = $1`

	GET_LOGS = `SELECT * FROM logs LIMIT 20;`

	GET_LOGFILES = `SELECT * FROM logfiles;`

	GET_LOGFILE = `SELECT * FROM logfiles WHERE logfileName LIKE $1;`

	GET_LOGFILE_NAMES = `SELECT id, logfileName FROM logfiles ORDER BY logfileName`

	GET_LOGS_SERVICE_INFO = `SELECT MIN(id), serviceInfo FROM %v GROUP BY serviceInfo;`

	GET_LOGS_SESSIONS = `SELECT * FROM %v;`

	GET_LOGS_BY_SESSION = `SELECT * FROM %v WHERE sessionId LIKE $1;`

	GET_ERRORS_BY_SESSION = `SELECT logId, sessionId, date, serviceInfo, description FROM %v WHERE sessionId LIKE $1;`

	GET_LOGS_BY_SESSION_FTS = `SELECT * FROM %v WHERE sessionId MATCH $1;`

	GET_ERRORS_BY_SESSION_LIMIT = `SELECT logId, sessionId, date, serviceInfo, description FROM %v WHERE sessionId LIKE $1 LIMIT $2 OFFSET $3;`

	GET_ERRORS_COUNT_BY_SESSION_LIMIT = `SELECT COUNT (*) as count FROM %v WHERE sessionId LIKE $1`

	GET_LOGS_BY_SESSION_LIMIT_FTS = `SELECT * FROM %v WHERE sessionId MATCH $1 LIMIT $2 OFFSET $3;`

	GET_LOGS_COUNT_BY_SESSION_LIMIT_FTS = `SELECT COUNT(*) as count FROM %v WHERE sessionId MATCH $1`

	GET_LOGS_BY_DATE_FTS = `SELECT * FROM %v WHERE date >= $1 AND date <= $2;`

	GET_ERRORS_BY_DATE = `SELECT logId, sessionId, date, serviceInfo, description FROM %v WHERE date >= $1 AND date <= $2;`

	GET_LOGS_BY_DATE_LIMIT_FTS = `SELECT * FROM %v WHERE date >= $1 AND date <= $2 LIMIT $3 OFFSET $4;`

	GET_ERRORS_BY_DATE_LIMIT = `SELECT logId, sessionId, date, serviceInfo, description FROM %v WHERE date >= $1 AND date <= $2 LIMIT $3 OFFSET $4;`

	GET_LOGS_COUNT_BY_DATE_LIMIT_FTS = `SELECT COUNT(*) as count FROM %v WHERE date >= $1 AND date <= $2;`

	GET_ERRORS_COUNT_BY_DATE_LIMIT = `SELECT COUNT(*) as count FROM %v WHERE date >= $1 AND date <= $2;`

	GET_ERRORS_BY_SESSION_FTS = `SELECT logId, sessionId, date, serviceInfo, description FROM %v WHERE sessionId MATCH $1;`

	GET_LOGS_BY_DESCRIPTION = `SELECT * FROM logs WHERE sessionId LIKE $1;`

	INSERT_LOG = "INSERT INTO %v (sessionId, date, serviceInfo, description) VALUES "

	INSERT_LOG_FILE = `INSERT INTO logfiles (logfileName, fileLength, lastSessionDate, previousDate, sessionCount) VALUES ($1, $2, $3, $4, $5);`

	UPDATE_LOG_FILE = `UPDATE logfiles SET logfileName = $1, fileLength = $2, lastSessionDate = $3, previousDate = $4, sessionCount = $5 WHERE logfileName = $1;`

	TRIGGER_FTS_INSERT = `CREATE TRIGGER %v AFTER INSERT ON %v
    	BEGIN
        	INSERT INTO %v (id, sessionId, date, serviceInfo, description)
        	VALUES (new.id, new.sessionId, new.date, new.serviceInfo, new.description);
    	END;`

	TRIGGER_CONDITIONAL_ERROR_INSERT1 = `CREATE TRIGGER %v AFTER INSERT ON %v
		WHEN %v LIKE '%%error%%'
    	BEGIN
        	INSERT INTO %v (logId, sessionId, date, serviceInfo, description)
        	VALUES (new.id, new.sessionId, new.date, new.serviceInfo, new.description);
    	END;`

	TRIGGER_CONDITIONAL_ERROR_INSERT2 = `CREATE TRIGGER %v AFTER INSERT ON %v
		WHEN %v LIKE '%%ошибк%%'
    	BEGIN
       	INSERT INTO %v (logId, sessionId, date, serviceInfo, description)
       	VALUES (new.id, new.sessionId, new.date, new.serviceInfo, new.description);
    	END;`

	FIND_LOG = `SELECT * FROM %[1]v WHERE %[1]v MATCH `

	FIND_LOG_COUNT_LIMIT = `SELECT COUNT (*) as count FROM %[1]v WHERE %[1]v MATCH `

	INSERT_LOG_SESSIONS = `INSERT INTO %v (dates) VALUES `

	INSERT_LOG_SESSION = `INSERT INTO %v (dates) VALUES ($1)`

	UPDATE_LOG_SESSION = `UPDATE %v SET dates = $1 WHERE id = $2`
)

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

func CreateLogDatabase(filename string) error {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	logsDbName := filename + "_logs"
	errorsDbName := filename + "_errors"

	_, err = db.Exec(fmt.Sprintf(CREATE_TABLE_LOGS, logsDbName))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = db.Exec(fmt.Sprintf(CREATE_TABLE_ERRORS, errorsDbName))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = db.Exec(fmt.Sprintf(CREATE_TABLE_SESSION_DATES, filename+"_sessions"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = db.Exec(fmt.Sprintf(CREATE_TABLE_LOGS_FTS5, filename+"_fts"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = db.Exec(fmt.Sprintf(TRIGGER_FTS_INSERT, filename+"_ai", logsDbName, filename+"_fts"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	var query string

	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT1, filename+"_errorInsert1", logsDbName, "new.serviceInfo", errorsDbName)
	_, err = db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT1, filename+"_errorInsert2", logsDbName, "new.description", errorsDbName)
	_, err = db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT2, filename+"_errorInsert3", logsDbName, "new.serviceInfo", errorsDbName)
	_, err = db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT2, filename+"_errorInsert4", logsDbName, "new.description", errorsDbName)
	_, err = db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}

	return nil
}

func Register(login string, password []byte, ip string) error {

	if len(password) == 0 {
		return erx.NewError(604, "Invalid password")
	}

	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	var user string
	err = db.QueryRow(GET_USER, login).Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}
	if err == nil {
		_, err = db.Exec(WRITE_STAT, 2, GetTime(), ip, "User already registered")
		if err != nil {
			return erx.New(err)
		}
		return erx.NewError(603, "User already registered")
	}

	salt, err := GenerateSalt()
	if err != nil {
		return erx.New(err)
	}
	_, err = db.Exec(WRITE_USER, login, PasswordHash(password, salt), salt, GetTime(), false, 5, 0, false)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(WRITE_STAT, 1, GetTime(), ip, "User added")

	return err
}

func Check(login string, password []byte, ip string) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	var user models.User
	err = db.QueryRow(GET_USER, login).Scan(&user.Login)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}

	if err == sql.ErrNoRows {
		_, err = db.Exec(WRITE_STAT, 3, GetTime(), ip, "User not found")
		return erx.NewError(604, "Invalid password")
	}

	err = db.QueryRow(GET_USER_PASSWORD, login).Scan(&user.Password, &user.Salt, &user.Blocked, &user.Try_count)
	if err != nil {
		return erx.New(err)
	}

	if user.Blocked {
		fmt.Println("User is blocked")
		_, err = db.Exec(WRITE_STAT, 5, GetTime(), ip, "Authentication attempt from a blocked user")
		return erx.NewError(604, "Invalid password")
	}

	passStr := hex.EncodeToString(user.Password)
	checkPass := PasswordHash(password, user.Salt)
	if passStr == hex.EncodeToString(checkPass) {
		fmt.Println("Password is correct")
		_, err = db.Exec(WRITE_TRY_COUNT, 5, login)
		if err != nil {
			return erx.New(err)
		}
		return nil
	}

	fmt.Println("Password is incorrect")
	user.Try_count--
	_, err = db.Exec(WRITE_STAT, 4, GetTime(), ip, "Invalid password")
	if err != nil {
		return erx.New(err)
	}
	_, err = db.Exec(WRITE_TRY_COUNT, user.Try_count, login)
	if err != nil {
		return erx.New(err)
	}

	if user.Try_count == 0 {
		fmt.Println("User is blocked")
		_, err = db.Exec(BLOCK_USER, GetTime(), login)
		if err != nil {
			return erx.New(err)
		}

		_, err = db.Exec(WRITE_STAT, 6, GetTime(), ip, "User blocked")
		if err != nil {
			return erx.New(err)
		}
	}

	return erx.NewError(604, "Invalid password")
}

func Block(login string, ip string) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	var blocked bool
	err = db.QueryRow(GET_BLOCKED, login).Scan(&blocked)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}

	if err == sql.ErrNoRows {
		fmt.Println("User not found")
		_, err := db.Exec(WRITE_STAT, 3, GetTime(), ip, "User not found")
		if err != nil {
			return erx.New(err)
		}
	}

	if blocked == true {
		return erx.NewError(601, "User is already blocked")
	}

	fmt.Println("User is blocked successfully")
	_, err = db.Exec(BLOCK_USER, GetTime(), login)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(WRITE_STAT, 6, GetTime(), ip, "User blocked")
	if err != nil {
		return erx.New(err)
	}

	return err
}

func Unblock(login string, ip string) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()
	var blocked bool
	err = db.QueryRow(GET_BLOCKED, login).Scan(&blocked)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}

	if err == sql.ErrNoRows {
		fmt.Println("User not found")
		_, err := db.Exec(WRITE_STAT, 3, GetTime(), ip, "User not found")
		if err != nil {
			return erx.New(err)
		}
	}

	if blocked == false {
		return erx.NewError(602, "User is not blocked")
	}

	fmt.Println("User is unblocked successfully")
	_, err = db.Exec(UNBLOCK_USER, login)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(WRITE_STAT, 7, GetTime(), ip, "User unblocked")
	if err != nil {
		return erx.New(err)
	}

	return err
}

func WriteRefreshToken(login string, token string, ua string, fingerprint string, ip string, daysUntilExpire int) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	_, err = db.Exec(WRITE_REFRESH_TOKEN, login, token, ua, fingerprint, ip, GetExpTime(daysUntilExpire), GetTime())
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(WRITE_STAT, 8, GetTime(), ip, "User auth success")
	if err != nil {
		return erx.New(err)
	}

	return err
}

func CheckRefreshToken(token string, fingerprint string, ip string) (string, int, error) {
	var login string
	var daysUntilExpire int
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return login, daysUntilExpire, erx.New(err)
	}
	defer db.Close()

	var refSes models.RefreshSession
	err = db.QueryRow(GET_REFRESH_SESSION, token).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.ExpiresIn, &refSes.CreatedAt)
	if err != nil && err != sql.ErrNoRows {
		return login, daysUntilExpire, erx.New(err)
	}

	if err == sql.ErrNoRows {
		tools.LogWarn("Session not found, Suspicious auth attempt")
		_, err := db.Exec(WRITE_STAT, 9, GetTime(), ip, "Suspicious auth attempt")
		if err != nil {
			return login, daysUntilExpire, erx.New(err)
		}
		//DeleteAllSessions()
		return login, daysUntilExpire, erx.New(sql.ErrNoRows)
	}

	_, err = db.Exec(DELETE_TOKEN, token)
	if err != nil {
		return login, daysUntilExpire, erx.New(err)
	}

	if refSes.ExpiresIn < GetTime() {
		tools.LogWarn("Session expired")
		_, err := db.Exec(WRITE_STAT, 10, GetTime(), ip, "Session expired")
		if err != nil {
			return login, daysUntilExpire, erx.New(err)
		}
		return login, daysUntilExpire, erx.NewError(608, "Session expired")
	}

	if fingerprint != refSes.Fingerprint {
		tools.LogWarn("Device not found, Suspicious auth attempt")
		_, err := db.Exec(WRITE_STAT, 9, GetTime(), ip, "Suspicious auth attempt")
		if err != nil {
			return login, daysUntilExpire, erx.New(err)
		}
		return login, daysUntilExpire, erx.NewError(615, "Suspicious device")
	}

	expiresIn := time.Unix(refSes.ExpiresIn, 0)
	createdAt := time.Unix(refSes.CreatedAt, 0)

	daysUntilExpire = int(expiresIn.Sub(createdAt).Hours()) / 24
	login = refSes.Login
	return login, daysUntilExpire, err
}

func DeleteSession(token string, ip string) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	var refSes models.RefreshSession
	err = db.QueryRow(FIND_REFRESH_SESSION, token).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.Ip)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}

	if err == sql.ErrNoRows {
		return nil
	}

	_, err = db.Exec(DELETE_TOKEN, refSes.Token)
	if err != nil {
		return erx.New(err)
	}
	_, err = db.Exec(WRITE_STAT, 12, GetTime(), ip, "Session deleted after user request")
	if err != nil {
		return erx.New(err)
	}

	return nil

}

func TryDeleteOldSession(fingerprint string, ip string) error {

	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	var refSes models.RefreshSession
	err = db.QueryRow(FIND_REFRESH_TOKEN, fingerprint).Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.Ip)
	if err != nil && err != sql.ErrNoRows {
		return erx.New(err)
	}

	if err == sql.ErrNoRows {
		return nil
	}

	_, err = db.Exec(DELETE_TOKEN, refSes.Token)
	if err != nil {
		return erx.New(err)
	}
	_, err = db.Exec(WRITE_STAT, 12, GetTime(), ip, "Session deleted after new login attempt")
	if err != nil {
		return erx.New(err)
	}

	return nil
}

// Function checks count of active Refresh Session. If there are already 3 sessions, func deletes the oldest session
func CheckSessionsCount(login string, ip string) error {
	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	rows, err := db.Query(GET_ALL_REFRESH_SESSIONS_SORTED, login)
	if err != nil {
		return erx.New(err)
	}
	defer rows.Close()

	var refSession models.RefreshSession
	var refSlice []models.RefreshSession

	for rows.Next() {
		err = rows.Scan(&refSession.Token, &refSession.Fingerprint)
		if err != nil {
			return erx.New(err)
		}
		refSlice = append(refSlice, refSession)
	}

	if len(refSlice) >= 3 {
		err = DeleteSession(refSlice[0].Token, ip)
		if err != nil {
			return erx.New(err)
		}
	}

	return nil

}

func DeleteAllSessions(login string, fingerprint string, ip string) error {

	db, err := sql.Open("sqlite3", getUsersDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	rows, err := db.Query(GET_ALL_REFRESH_SESSIONS, login)
	if err != nil {
		return erx.New(err)
	}
	defer rows.Close()

	var refSes models.RefreshSession
	var tkns []string
	var sessionValid = false
	for rows.Next() {
		err = rows.Scan(&refSes.Login, &refSes.Token, &refSes.Fingerprint, &refSes.Ip)
		if err != nil {
			return erx.New(err)
		}
		if refSes.Fingerprint == fingerprint {
			sessionValid = true
			continue
		}
		tkns = append(tkns, refSes.Token)
	}
	if tkns == nil {
		return erx.NewError(614, "No sessions found")
	}
	if sessionValid {
		for _, v := range tkns {

			_, err := db.Exec(DELETE_TOKEN, v)
			if err != nil {
				return erx.New(err)
			}
			_, err = db.Exec(WRITE_STAT, 12, GetTime(), ip, "Session deleted after Exit Everywhere action")
			if err != nil {
				return erx.New(err)
			}
		}
		return nil
	}

	return erx.NewError(617, "No valid session found")
}

func InsertLogs(slice []models.LogStruct, filename string) (int64, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return 0, erx.New(err)
	}
	defer db.Close()

	var values []interface{}
	var builder strings.Builder
	query := fmt.Sprintf(INSERT_LOG, filename)
	builder.WriteString(query)
	var rowCount int64

	for i, v := range slice {
		if i%249 == 0 && i != 0 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := db.Prepare(reqString)
			if err != nil {
				tools.LogErr(err)
				return rowCount, err
			}
			res, err := request.Exec(values...)
			if err != nil {
				tools.LogErr(err)
				return rowCount, err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				tools.LogErr(erx.New(err))
				return rowCount, err
			}
			rowCount += rows
			builder.Reset()
			builder.WriteString(query)
			values = nil
			continue
		}

		builder.WriteString("(?, ?, ?, ?),")
		values = append(values, v.SessionId, v.Date, v.ServiceInfo, v.Description)

		if i == len(slice)-1 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := db.Prepare(reqString)
			if err != nil {
				tools.LogErr(erx.New(err))
				return rowCount, err
			}
			res, err := request.Exec(values...)
			if err != nil {
				tools.LogErr(erx.New(err))
				return rowCount, err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				tools.LogErr(erx.New(err))
				return rowCount, err
			}
			rowCount += rows

		}

	}
	return rowCount, nil
}

func InsertLogFileInfo(fileStruct models.LogFileStruct) error {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	_, err = db.Exec(INSERT_LOG_FILE, fileStruct.LogfileName, fileStruct.FileLength, fileStruct.LastSessionDate, fileStruct.PreviousDate, fileStruct.SessionCount)
	if err != nil {
		return erx.New(err)
	}

	return err
}

func UpdateLogFileInfo(fileStruct models.LogFileStruct) error {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	_, err = db.Exec(UPDATE_LOG_FILE, fileStruct.LogfileName, fileStruct.FileLength, fileStruct.LastSessionDate, fileStruct.PreviousDate, fileStruct.SessionCount)
	if err != nil {
		return erx.New(err)
	}

	return err
}

func InsertLogSessions(slice []models.LogSessionStruct, filename string) error {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()

	var values []interface{}
	var builder strings.Builder
	builder.WriteString(INSERT_LOG_SESSIONS)

	for i, v := range slice {
		if i%499 == 0 && i != 0 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := db.Prepare(reqString)
			if err != nil {
				tools.LogErr(err)
				return err
			}
			_, err = request.Exec(values...)
			if err != nil {
				tools.LogErr(err)
				return err
			}
			builder.Reset()
			query := fmt.Sprintf(INSERT_LOG, filename)
			builder.WriteString(query)
			values = nil
			continue
		}

		builder.WriteString("(?),")
		values = append(values, v.Dates)

		if i == len(slice)-1 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := db.Prepare(reqString)
			if err != nil {
				tools.LogErr(err)
				return err
			}
			res, err := request.Exec(values...)
			if err != nil {
				tools.LogErr(err)
				return err
			}
			_, err = res.RowsAffected()
			if err != nil {
				tools.LogErr(err)
				return err
			}
		}

	}

	return nil
}

func InsertLogSession(session models.LogSessionStruct, filename string) error {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return erx.New(err)
	}
	defer db.Close()
	query := fmt.Sprintf(UPDATE_LOG_SESSION, filename+"_sessions")
	res, err := db.Exec(query, session.Dates, session.Id)
	if err != nil {
		return erx.New(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}

	if rows == 0 {
		query := fmt.Sprintf(INSERT_LOG_SESSION, filename+"_sessions")
		_, err = db.Exec(query, session.Dates)
	}

	return err
}

func GetLogsBySession(sessionId int, filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if sessionId == 0 {
		sessionCount := os.Getenv(filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		sessionId = countInt
	}
	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_SESSION_FTS, filename+"_fts")
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogStruct
	var logSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
		logSlice = append(logSlice, logRow)
	}

	if len(logSlice) == 0 {
		noRows := logsNotFoundStruct
		logSlice = append(logSlice, noRows)
	}

	s := fmt.Sprintf("Sending session #%v", sessionId)

	tools.LogInfo(s)

	return logSlice, nil
}

func GetErrorsBySession(sessionId int, filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if sessionId == 0 {
		sessionCount := os.Getenv(filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		sessionId = countInt
	}

	query := fmt.Sprintf(GET_ERRORS_BY_SESSION, filename+"_errors")
	rows, err := db.Query(query, sessionId)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()
	var errorRow models.LogStruct
	var errorSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&errorRow.Id, &errorRow.SessionId, &errorRow.Date, &errorRow.ServiceInfo, &errorRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		errorRow.DateUtc = tools.FormatUnixToUTC(errorRow.Date)
		errorSlice = append(errorSlice, errorRow)
	}

	if len(errorSlice) == 0 {
		noRows := logsNotFoundStruct
		errorSlice = append(errorSlice, noRows)
	}

	return errorSlice, nil
}

func GetLogsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_DATE_FTS, filename+"_fts")
	rows, err := db.Query(query, startDate.Unix(), endDate.Unix())
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogStruct
	var logSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
		logSlice = append(logSlice, logRow)
	}

	if len(logSlice) == 0 {
		noRows := logsNotFoundStruct
		logSlice = append(logSlice, noRows)
	}

	s := fmt.Sprintf("Sending session between dates %v and %v", startDate.UTC(), endDate.UTC())

	tools.LogInfo(s)

	return logSlice, nil
}

func GetErrorsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	query := fmt.Sprintf(GET_ERRORS_BY_DATE, filename+"_errors")
	rows, err := db.Query(query, startDate.Unix(), endDate.Unix())
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()
	var errorRow models.LogStruct
	var errorSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&errorRow.Id, &errorRow.SessionId, &errorRow.Date, &errorRow.ServiceInfo, &errorRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		errorRow.DateUtc = tools.FormatUnixToUTC(errorRow.Date)
		errorSlice = append(errorSlice, errorRow)
	}

	if len(errorSlice) == 0 {
		noRows := logsNotFoundStruct
		errorSlice = append(errorSlice, noRows)
	}

	return errorSlice, nil
}

func GetLogsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {
	var rowsOffset = limit
	var count int
	offset = offset - rowsOffset
	logSlices := [3][]models.LogStruct{}
	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		err = erx.New(err)
		return logSlices, count, err
	}

	defer db.Close()

	if filename == "" {
		err = erx.NewError(0, "Error: filename is not specified")
		return logSlices, count, err
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	query := fmt.Sprintf(GET_LOGS_BY_DATE_LIMIT_FTS, filename+"_fts")

	var logRow models.LogStruct
	for i, slice := range logSlices {
		if offset < 0 {
			offset = offset + rowsOffset
			continue
		}

		rows, err := db.Query(query, startDate.Unix(), endDate.Unix(), limit, offset)
		if err != nil {
			err = erx.New(err)
			rows.Close()
			return logSlices, count, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
			if err != nil {
				tools.LogErr(err)
				err = erx.New(err)
				return logSlices, count, err
			}
			logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
			slice = append(slice, logRow)
		}
		logSlices[i] = slice
		if offset == 0 {
			logSlices[0] = slice
		}
		offset = offset + rowsOffset
	}

	//if len(logSlices[1]) == 0 {
	//	noRows := logsNotFoundStruct
	//	logSlices[1] = append(logSlices[1], noRows)
	//}

	if len(logSlices[1]) == 0 {
		logSlices[1] = logSlices[0]
	}

	if len(logSlices[1]) != 0 {
		query := fmt.Sprintf(GET_LOGS_COUNT_BY_DATE_LIMIT_FTS, filename+"_fts")
		err := db.QueryRow(query, startDate.Unix(), endDate.Unix()).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return logSlices, count, err
		}
	}

	if logSlices[len(logSlices)-1] == nil {
		logSlices[len(logSlices)-1] = logSlices[len(logSlices)-2]
	}

	s := fmt.Sprintf("Sending part of Logs slice between dates %v and %v", startDate.UTC(), endDate.UTC())

	tools.LogInfo(s)

	return logSlices, count, nil
}

func GetErrorsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {

	var rowsOffset = limit
	var count int
	offset = offset - rowsOffset
	errorSlices := [3][]models.LogStruct{}

	//var errorSlice []models.LogStruct
	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		err = erx.New(err)
		return errorSlices, count, err
	}
	defer db.Close()

	if filename == "" {
		err = erx.NewError(0, "Error: filename is not specified")
		return errorSlices, count, err
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	query := fmt.Sprintf(GET_ERRORS_BY_DATE_LIMIT, filename+"_errors")
	var errorRow models.LogStruct
	for i, slice := range errorSlices {
		if offset < 0 {
			offset = offset + rowsOffset
			continue
		}
		rows, err := db.Query(query, startDate.Unix(), endDate.Unix(), limit, offset)
		if err != nil {
			err = erx.New(err)
			rows.Close()
			return errorSlices, count, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&errorRow.Id, &errorRow.SessionId, &errorRow.Date, &errorRow.ServiceInfo, &errorRow.Description)
			if err != nil {
				tools.LogErr(err)
				err = erx.New(err)
				return errorSlices, count, err
			}
			errorRow.DateUtc = tools.FormatUnixToUTC(errorRow.Date)
			slice = append(slice, errorRow)
		}
		errorSlices[i] = slice
		if offset == 0 {
			errorSlices[0] = slice
		}
		offset = offset + rowsOffset
	}
	//rows, err := db.Query(query, startDate.Unix(), endDate.Unix(), limit, offset)
	//if err != nil{
	//	err = erx.New(err)
	//	return errorSlice, count, err
	//}
	//defer rows.Close()
	//var errorRow models.LogStruct

	//for rows.Next() {
	//	err = rows.Scan(&errorRow.Id, &errorRow.SessionId, &errorRow.Date, &errorRow.ServiceInfo, &errorRow.Description)
	//	if err != nil {
	//		tools.LogErr(err)
	//		err = erx.New(err)
	//		return errorSlice, count, err
	//	}
	//	errorRow.DateUtc = tools.FormatUnixToUTC(errorRow.Date)
	//	errorSlice = append(errorSlice, errorRow)
	//}

	if len(errorSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		errorSlices[1] = append(errorSlices[1], noRows)
	}

	if len(errorSlices[1]) != 0 {
		query := fmt.Sprintf(GET_ERRORS_COUNT_BY_DATE_LIMIT, filename+"_errors")
		err := db.QueryRow(query, startDate.Unix(), endDate.Unix()).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return errorSlices, count, err
		}
	}

	return errorSlices, count, nil
}

func GetLogsBySessionWithLimit(sessionId int, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {

	var rowsOffset = limit
	var count int
	var logRow models.LogStruct
	offset = offset - rowsOffset

	logSlices := [3][]models.LogStruct{}

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		err = erx.New(err)
		return logSlices, count, err
	}
	defer db.Close()

	if filename == "" {
		err = erx.NewError(0, "Error: filename is not specified")
		return logSlices, count, err
	}

	if sessionId == 0 {
		sessionCount := os.Getenv(filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		sessionId = countInt
	}
	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_SESSION_LIMIT_FTS, filename+"_fts")

	for i, slice := range logSlices {
		if offset < 0 {
			offset = offset + rowsOffset
			continue
		}

		rows, err := db.Query(query, sessionId, limit, offset)
		if err != nil {
			err = erx.New(err)
			rows.Close()
			return logSlices, count, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
			if err != nil {
				tools.LogErr(err)
				err = erx.New(err)
				return logSlices, count, err
			}
			logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
			slice = append(slice, logRow)
		}
		logSlices[i] = slice
		if offset == 0 {
			logSlices[0] = slice
		}
		offset = offset + rowsOffset
	}

	if len(logSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		logSlices[1] = append(logSlices[1], noRows)
	}

	if len(logSlices[1]) != 0 {
		query := fmt.Sprintf(GET_LOGS_COUNT_BY_SESSION_LIMIT_FTS, filename+"_fts")
		err := db.QueryRow(query, sessionId).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return logSlices, count, err
		}
	}

	if logSlices[len(logSlices)-1] == nil {
		logSlices[len(logSlices)-1] = logSlices[len(logSlices)-2]
	}

	s := fmt.Sprintf("Sending session #%v", sessionId)

	tools.LogInfo(s)

	return logSlices, count, nil
}

func GetErrorsBySessionWithLimit(sessionId int, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {

	var rowsOffset = limit
	var count int
	var errorRow models.LogStruct
	offset = offset - rowsOffset

	errorSlices := [3][]models.LogStruct{}

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		err = erx.New(err)
		return errorSlices, count, err
	}
	defer db.Close()

	if filename == "" {
		err = erx.NewError(0, "Error: filename is not specified")
		return errorSlices, count, err
	}

	if sessionId == 0 {
		sessionCount := os.Getenv(filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		sessionId = countInt
	}

	query := fmt.Sprintf(GET_ERRORS_BY_SESSION_LIMIT, filename+"_errors")
	//rows, err := db.Query(query, sessionId, limit, offset)

	for i, slice := range errorSlices {
		if offset < 0 {
			offset = offset + rowsOffset
			continue
		}

		rows, err := db.Query(query, sessionId, limit, offset)
		if err != nil {
			err = erx.New(err)
			rows.Close()
			return errorSlices, count, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&errorRow.Id, &errorRow.SessionId, &errorRow.Date, &errorRow.ServiceInfo, &errorRow.Description)
			if err != nil {
				tools.LogErr(err)
				err = erx.New(err)
				return errorSlices, count, err
			}
			errorRow.DateUtc = tools.FormatUnixToUTC(errorRow.Date)
			slice = append(slice, errorRow)
		}
		errorSlices[i] = slice
		if offset == 0 {
			errorSlices[0] = slice
		}
		offset = offset + rowsOffset
	}

	if len(errorSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		errorSlices[1] = append(errorSlices[1], noRows)
	}

	if len(errorSlices[1]) != 0 {
		query := fmt.Sprintf(GET_ERRORS_COUNT_BY_SESSION_LIMIT, filename+"_errors")
		err := db.QueryRow(query, sessionId).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return errorSlices, count, err
		}
	}

	if errorSlices[len(errorSlices)-1] == nil {
		errorSlices[len(errorSlices)-1] = errorSlices[len(errorSlices)-2]
	}

	return errorSlices, count, nil
}

func FindLogsWithLimit(searchText string, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {

	var rowsOffset = limit
	var count int
	var foundRow models.LogStruct
	offset = offset - rowsOffset

	foundSlices := [3][]models.LogStruct{}

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return foundSlices, count, erx.New(err)
	}
	defer db.Close()

	query := fmt.Sprintf(FIND_LOG, filename+"_fts")
	searchText = fmt.Sprintf("'\"%s\"*';", searchText)
	completeQuery := query + searchText

	for i, slice := range foundSlices {
		if offset < 0 {
			offset = offset + rowsOffset
			continue
		}
		rows, err := db.Query(completeQuery)
		if err != nil {
			rows.Close()
			return foundSlices, count, erx.New(err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&foundRow.Id, &foundRow.SessionId, &foundRow.Date, &foundRow.ServiceInfo, &foundRow.Description)
			if err != nil {
				tools.LogErr(err)
				return foundSlices, count, err
			}
			foundRow.DateUtc = tools.FormatUnixToUTC(foundRow.Date)
			slice = append(slice, foundRow)
		}
		foundSlices[i] = slice
		if offset == 0 {
			foundSlices[0] = slice
		}
		offset = offset + rowsOffset
	}

	if len(foundSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		foundSlices[1] = append(foundSlices[1], noRows)
	}

	if len(foundSlices[1]) != 0 {
		query := fmt.Sprintf(FIND_LOG_COUNT_LIMIT, filename+"_fts")
		completeQuery = query + searchText
		//query := fmt.Sprintf(GET_ERRORS_COUNT_BY_SESSION_LIMIT, filename + "_errors")
		err := db.QueryRow(completeQuery).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return foundSlices, count, err
		}
	}

	if foundSlices[len(foundSlices)-1] == nil {
		foundSlices[len(foundSlices)-1] = foundSlices[len(foundSlices)-2]
	}

	//s := fmt.Sprintf("Found: \n", logSlice)

	//tools.LogInfo(s)

	return foundSlices, count, nil
}

func FindLogs(searchText string, filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	query := fmt.Sprintf(FIND_LOG, filename+"_fts")
	searchText = fmt.Sprintf("'\"%s\"*';", searchText)
	completeQuery := query + searchText
	rows, err := db.Query(completeQuery)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogStruct
	var logSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
		logSlice = append(logSlice, logRow)
	}

	if len(logSlice) == 0 {
		noRows := logsNotFoundStruct
		logSlice = append(logSlice, noRows)
	}

	//s := fmt.Sprintf("Found: \n", logSlice)

	//tools.LogInfo(s)

	return logSlice, nil
}

func GetLogById(id int, filename string) (models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return models.LogStruct{}, erx.New(err)
	}
	defer db.Close()

	var logRow models.LogStruct

	query := fmt.Sprintf(GET_LOG_BY_ID, filename+"_logs")
	err = db.QueryRow(query, id).Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
	if err != nil && err != sql.ErrNoRows {
		tools.LogErr(erx.New(err))
		return models.LogStruct{}, erx.New(err)
	}
	if err == sql.ErrNoRows {
		logRow = logsNotFoundStruct
		tools.LogErr(erx.New(err))
		return logRow, erx.New(err)
	}
	logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
	s := fmt.Sprintf("Found:%v", logRow)
	tools.LogInfo(s)

	return logRow, nil
}

func GetLogs() ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	rows, err := db.Query(GET_LOGS)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogStruct
	var logSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logRow.DateUtc = tools.FormatUnixToUTC(logRow.Date)
		logSlice = append(logSlice, logRow)
	}

	return logSlice, nil
}

func GetLogsSessions(filename string) ([]models.LogSessionStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	query := fmt.Sprintf(GET_LOGS_SESSIONS, filename+"_sessions")
	rows, err := db.Query(query)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logSessionRow models.LogSessionStruct
	var logSessionSlice []models.LogSessionStruct
	for rows.Next() {
		err = rows.Scan(&logSessionRow.Id, &logSessionRow.Dates)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logSessionSlice = append(logSessionSlice, logSessionRow)
	}

	return logSessionSlice, nil
}

func GetLogsServiceInfo(filename string) ([]models.LogStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	query := fmt.Sprintf(GET_LOGS_SERVICE_INFO, filename+"_logs")
	rows, err := db.Query(query)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logSessionRow models.LogStruct
	var logSessionSlice []models.LogStruct
	for rows.Next() {
		err = rows.Scan(&logSessionRow.Id, &logSessionRow.ServiceInfo)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logSessionSlice = append(logSessionSlice, logSessionRow)
	}

	return logSessionSlice, nil
}

func GetLogsFileNames() ([]models.LogFilenameStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	rows, err := db.Query(GET_LOGFILE_NAMES)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogFilenameStruct
	var logSlice []models.LogFilenameStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.LogfileName)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logSlice = append(logSlice, logRow)
	}

	return logSlice, nil
}

func GetLogsFilesInfo() ([]models.LogFileStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return nil, erx.New(err)
	}
	defer db.Close()

	rows, err := db.Query(GET_LOGFILES)
	if err != nil {
		return nil, erx.New(err)
	}
	defer rows.Close()

	var logRow models.LogFileStruct
	var logSlice []models.LogFileStruct
	for rows.Next() {
		err = rows.Scan(&logRow.Id, &logRow.LogfileName, &logRow.FileLength, &logRow.LastSessionDate, &logRow.PreviousDate, &logRow.SessionCount)
		if err != nil {
			tools.LogErr(err)
			return nil, err
		}
		logSlice = append(logSlice, logRow)
	}

	return logSlice, nil
}

func GetLogsFileInfo(filename string) (models.LogFileStruct, error) {

	db, err := sql.Open("sqlite3", getLogsDbPath())
	if err != nil {
		return models.LogFileStruct{}, erx.New(err)
	}
	defer db.Close()

	var logRow models.LogFileStruct
	err = db.QueryRow(GET_LOGFILE, filename).Scan(&logRow.Id, &logRow.LogfileName, &logRow.FileLength, &logRow.LastSessionDate, &logRow.PreviousDate, &logRow.SessionCount)
	if err != nil {
		return models.LogFileStruct{}, err
	}

	return logRow, nil
}
