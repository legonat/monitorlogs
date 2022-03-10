package db

import (
	"database/sql"
	"monitorlogs/internal/models"
	"time"
)

type UsersDb interface {
	Register(inputs models.RegisterInputs) error
	Check(inputs models.LoginInputs) error
	Block(inputs models.BlockInputs) error
	Unblock(inputs models.BlockInputs) error
	WriteRefreshToken(inputs models.RefreshSession, daysUntilExpire int) error
	CheckRefreshToken(inputs models.RefreshSession) (string, int, error)
	DeleteSession(token string, ip string) error
	TryDeleteOldSession(fingerprint string, ip string) error
	CheckSessionsCount(login string, ip string) error
	DeleteAllSessions(inputs models.RefreshSession) error
}

type LogsDb interface {
	CreateLogDatabase(filename string) error
	InsertLogs(slice []models.LogStruct, filename string) (int64, error)
	InsertLogFileInfo(fileStruct models.LogFileStruct) error
	UpdateLogFileInfo(fileStruct models.LogFileStruct) error
	InsertLogSessions(slice []models.LogSessionStruct, filename string) error
	InsertLogSession(session models.LogSessionStruct, filename string) error
	GetLogsBySession(inputs models.GetLogsBySessionStruct) ([]models.LogStruct, error)
	GetErrorsBySession(inputs models.GetLogsBySessionStruct) ([]models.LogStruct, error)
	GetLogsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error)
	GetErrorsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error)
	GetLogsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error)
	GetErrorsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error)
	GetLogsBySessionWithLimit(inputs models.GetLogsBySessionWithLimitStruct) ([3][]models.LogStruct, int, error)
	GetErrorsBySessionWithLimit(inputs models.GetLogsBySessionWithLimitStruct) ([3][]models.LogStruct, int, error)
	FindLogsWithLimit(inputs models.FindLogsStructWithLimit) ([3][]models.LogStruct, int, error)
	FindLogs(inputs models.FindLogsStruct) ([]models.LogStruct, error)
	GetLogById(id int, filename string) (models.LogStruct, error)
	GetLogs() ([]models.LogStruct, error)
	GetLogsSessions(filename string) ([]models.LogSessionStruct, error)
	GetLogsServiceInfo(filename string) ([]models.LogStruct, error)
	GetLogsFileNames() ([]models.LogFilenameStruct, error)
	GetLogsFilesInfo() ([]models.LogFileStruct, error)
	GetLogsFileInfo(filename string) (models.LogFileStruct, error)
	ReadFolder(folderPath string)
	ReadCycle(duration string, folderPath string)
	Read(fullFilename string) error
}

type Repository struct {
	LogsDb
	UsersDb
}

func NewRepository(logsDb, usersDb *sql.DB) *Repository {
	return &Repository{
		LogsDb : NewLogsDbSqlite(logsDb),
		UsersDb: NewUsersDbSqlite(usersDb),
	}
}