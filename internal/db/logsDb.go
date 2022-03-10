package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"
	"strconv"
	"strings"
	"time"
)

const(

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


	CREATE_TABLE_SESSION_DATES = `
	CREATE TABLE %v (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	dates INTEGER NOT NULL);
	`

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

type LogsDbSqlite struct {
	db *sql.DB
}

func NewLogsDbSqlite(db *sql.DB) *LogsDbSqlite {
	return &LogsDbSqlite{db: db}
}

func (r *LogsDbSqlite) CreateLogDatabase(filename string) error {


	logsDbName := filename + "_logs"
	errorsDbName := filename + "_errors"

	_, err := r.db.Exec(fmt.Sprintf(CREATE_TABLE_LOGS, logsDbName))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = r.db.Exec(fmt.Sprintf(CREATE_TABLE_ERRORS, errorsDbName))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = r.db.Exec(fmt.Sprintf(CREATE_TABLE_SESSION_DATES, filename+"_sessions"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = r.db.Exec(fmt.Sprintf(CREATE_TABLE_LOGS_FTS5, filename+"_fts"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	_, err = r.db.Exec(fmt.Sprintf(TRIGGER_FTS_INSERT, filename+"_ai", logsDbName, filename+"_fts"))
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	var query string

	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT1, filename+"_errorInsert1", logsDbName, "new.serviceInfo", errorsDbName)
	_, err = r.db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT1, filename+"_errorInsert2", logsDbName, "new.description", errorsDbName)
	_, err = r.db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT2, filename+"_errorInsert3", logsDbName, "new.serviceInfo", errorsDbName)
	_, err = r.db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}
	query = fmt.Sprintf(TRIGGER_CONDITIONAL_ERROR_INSERT2, filename+"_errorInsert4", logsDbName, "new.description", errorsDbName)
	_, err = r.db.Exec(query)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}

	return nil
}

func (r *LogsDbSqlite) InsertLogs(slice []models.LogStruct, filename string) (int64, error) {

	var values []interface{}
	var builder strings.Builder
	query := fmt.Sprintf(INSERT_LOG, filename)
	builder.WriteString(query)
	var rowCount int64

	for i, v := range slice {
		if i%249 == 0 && i != 0 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := r.db.Prepare(reqString)
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
			request, err := r.db.Prepare(reqString)
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

func (r *LogsDbSqlite) InsertLogFileInfo(fileStruct models.LogFileStruct) error {

	_, err := r.db.Exec(INSERT_LOG_FILE, fileStruct.LogfileName, fileStruct.FileLength, fileStruct.LastSessionDate, fileStruct.PreviousDate, fileStruct.SessionCount)
	if err != nil {
		return erx.New(err)
	}

	return err
}

func (r *LogsDbSqlite) UpdateLogFileInfo(fileStruct models.LogFileStruct) error {

	_, err := r.db.Exec(UPDATE_LOG_FILE, fileStruct.LogfileName, fileStruct.FileLength, fileStruct.LastSessionDate, fileStruct.PreviousDate, fileStruct.SessionCount)
	if err != nil {
		return erx.New(err)
	}

	return err
}

func (r *LogsDbSqlite) InsertLogSessions(slice []models.LogSessionStruct, filename string) error {


	var values []interface{}
	var builder strings.Builder
	builder.WriteString(INSERT_LOG_SESSIONS)

	for i, v := range slice {
		if i%499 == 0 && i != 0 {
			reqString := builder.String()
			reqString = reqString[:builder.Len()-1]
			request, err := r.db.Prepare(reqString)
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
			request, err := r.db.Prepare(reqString)
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

func (r *LogsDbSqlite) InsertLogSession(session models.LogSessionStruct, filename string) error {

	query := fmt.Sprintf(UPDATE_LOG_SESSION, filename+"_sessions")
	res, err := r.db.Exec(query, session.Dates, session.Id)
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
		_, err = r.db.Exec(query, session.Dates)
	}

	return err
}

func (r *LogsDbSqlite) GetLogsBySession(inputs models.GetLogsBySessionStruct) ([]models.LogStruct, error) {

	if inputs.Filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if inputs.SessionId == 0 {
		sessionCount := os.Getenv(inputs.Filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		inputs.SessionId = countInt
	}
	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_SESSION_FTS, inputs.Filename+"_fts")
	rows, err := r.db.Query(query, inputs.SessionId)
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

	s := fmt.Sprintf("Sending session #%v", inputs.SessionId)

	tools.LogInfo(s)

	return logSlice, nil
}

func (r *LogsDbSqlite) GetErrorsBySession(inputs models.GetLogsBySessionStruct) ([]models.LogStruct, error) {

	if inputs.Filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if inputs.SessionId == 0 {
		sessionCount := os.Getenv(inputs.Filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		inputs.SessionId = countInt
	}

	query := fmt.Sprintf(GET_ERRORS_BY_SESSION, inputs.Filename+"_errors")
	rows, err := r.db.Query(query, inputs.SessionId)
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

func (r *LogsDbSqlite) GetLogsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error) {

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_DATE_FTS, filename+"_fts")
	rows, err := r.db.Query(query, startDate.Unix(), endDate.Unix())
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

func (r *LogsDbSqlite) GetErrorsByDate(startDate time.Time, endDate time.Time, filename string) ([]models.LogStruct, error) {

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	if endDate.Unix() == 0 {
		startDate.AddDate(0, 0, 1)
	}

	query := fmt.Sprintf(GET_ERRORS_BY_DATE, filename+"_errors")
	rows, err := r.db.Query(query, startDate.Unix(), endDate.Unix())
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

func (r *LogsDbSqlite) GetLogsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {
	var rowsOffset = limit
	var count int
	offset = offset - rowsOffset
	logSlices := [3][]models.LogStruct{}

	if filename == "" {
		err := erx.NewError(0, "Error: filename is not specified")
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

		rows, err := r.db.Query(query, startDate.Unix(), endDate.Unix(), limit, offset)
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
		err := r.db.QueryRow(query, startDate.Unix(), endDate.Unix()).Scan(&count)
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

func (r *LogsDbSqlite) GetErrorsByDateWithLimit(startDate time.Time, endDate time.Time, filename string, limit int, offset int) ([3][]models.LogStruct, int, error) {

	var rowsOffset = limit
	var count int
	offset = offset - rowsOffset
	errorSlices := [3][]models.LogStruct{}


	if filename == "" {
		err := erx.NewError(0, "Error: filename is not specified")
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
		rows, err := r.db.Query(query, startDate.Unix(), endDate.Unix(), limit, offset)
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
		query := fmt.Sprintf(GET_ERRORS_COUNT_BY_DATE_LIMIT, filename+"_errors")
		err := r.db.QueryRow(query, startDate.Unix(), endDate.Unix()).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return errorSlices, count, err
		}
	}

	return errorSlices, count, nil
}

func (r *LogsDbSqlite) GetLogsBySessionWithLimit(inputs models.GetLogsBySessionWithLimitStruct) ([3][]models.LogStruct, int, error) {

	var rowsOffset = inputs.Limit
	var count int
	var logRow models.LogStruct
	inputs.Offset = inputs.Offset - rowsOffset

	logSlices := [3][]models.LogStruct{}

	if inputs.Filename == "" {
		err := erx.NewError(0, "Error: filename is not specified")
		return logSlices, count, err
	}

	if inputs.SessionId == 0 {
		sessionCount := os.Getenv(inputs.Filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		inputs.SessionId = countInt
	}
	//query := fmt.Sprintf(GET_LOGS_BY_SESSION, filename + "_logs")
	query := fmt.Sprintf(GET_LOGS_BY_SESSION_LIMIT_FTS, inputs.Filename+"_fts")

	for i, slice := range logSlices {
		if inputs.Offset < 0 {
			inputs.Offset = inputs.Offset + rowsOffset
			continue
		}

		rows, err := r.db.Query(query, inputs.SessionId, inputs.Limit, inputs.Offset)
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
		if inputs.Offset == 0 {
			logSlices[0] = slice
		}
		inputs.Offset = inputs.Offset + rowsOffset
	}

	if len(logSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		logSlices[1] = append(logSlices[1], noRows)
	}

	if len(logSlices[1]) != 0 {
		query := fmt.Sprintf(GET_LOGS_COUNT_BY_SESSION_LIMIT_FTS, inputs.Filename+"_fts")
		err := r.db.QueryRow(query, inputs.SessionId).Scan(&count)
		if err != nil {
			err = erx.New(err)
			return logSlices, count, err
		}
	}

	if logSlices[len(logSlices)-1] == nil {
		logSlices[len(logSlices)-1] = logSlices[len(logSlices)-2]
	}

	s := fmt.Sprintf("Sending session #%v", inputs.SessionId)

	tools.LogInfo(s)

	return logSlices, count, nil
}

func (r *LogsDbSqlite) GetErrorsBySessionWithLimit(inputs models.GetLogsBySessionWithLimitStruct) ([3][]models.LogStruct, int, error) {

	var rowsOffset = inputs.Limit
	var count int
	var errorRow models.LogStruct
	inputs.Offset = inputs.Offset - rowsOffset

	errorSlices := [3][]models.LogStruct{}

	if inputs.Filename == "" {
		err := erx.NewError(0, "Error: filename is not specified")
		return errorSlices, count, err
	}

	if inputs.SessionId == 0 {
		sessionCount := os.Getenv(inputs.Filename + "SessionCount")
		countInt, err := strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			countInt = 1
		}
		inputs.SessionId = countInt
	}

	query := fmt.Sprintf(GET_ERRORS_BY_SESSION_LIMIT, inputs.Filename+"_errors")
	//rows, err := r.db.Query(query, sessionId, limit, offset)

	for i, slice := range errorSlices {
		if inputs.Offset < 0 {
			inputs.Offset = inputs.Offset + rowsOffset
			continue
		}

		rows, err := r.db.Query(query, inputs.SessionId, inputs.Limit, inputs.Offset)
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
		if inputs.Offset == 0 {
			errorSlices[0] = slice
		}
		inputs.Offset = inputs.Offset + rowsOffset
	}

	if len(errorSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		errorSlices[1] = append(errorSlices[1], noRows)
	}

	if len(errorSlices[1]) != 0 {
		query := fmt.Sprintf(GET_ERRORS_COUNT_BY_SESSION_LIMIT, inputs.Filename+"_errors")
		err := r.db.QueryRow(query, inputs.SessionId).Scan(&count)
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

func (r *LogsDbSqlite) FindLogsWithLimit(inputs models.FindLogsStructWithLimit) ([3][]models.LogStruct, int, error) {

	var rowsOffset = inputs.Limit
	var count int
	var foundRow models.LogStruct
	inputs.Offset = inputs.Offset - rowsOffset

	foundSlices := [3][]models.LogStruct{}

	query := fmt.Sprintf(FIND_LOG, inputs.Filename+"_fts")
	inputs.SearchText = fmt.Sprintf("'\"%s\"*';", inputs.SearchText)
	completeQuery := query + inputs.SearchText

	for i, slice := range foundSlices {
		if inputs.Offset < 0 {
			inputs.Offset = inputs.Offset + rowsOffset
			continue
		}
		rows, err := r.db.Query(completeQuery)
		if err != nil {
			rows.Close()
			return foundSlices, count, erx.New(err)
		}

		for rows.Next() {
			err = rows.Scan(&foundRow.Id, &foundRow.SessionId, &foundRow.Date, &foundRow.ServiceInfo, &foundRow.Description)
			if err != nil {
				rows.Close()
				tools.LogErr(err)
				return foundSlices, count, err
			}
			foundRow.DateUtc = tools.FormatUnixToUTC(foundRow.Date)
			slice = append(slice, foundRow)
		}
		foundSlices[i] = slice
		if inputs.Offset == 0 {
			foundSlices[0] = slice
		}
		inputs.Offset = inputs.Offset + rowsOffset
		rows.Close()
	}

	if len(foundSlices[1]) == 0 {
		noRows := logsNotFoundStruct
		foundSlices[1] = append(foundSlices[1], noRows)
	}

	if len(foundSlices[1]) != 0 {
		query := fmt.Sprintf(FIND_LOG_COUNT_LIMIT, inputs.Filename+"_fts")
		completeQuery = query + inputs.SearchText
		//query := fmt.Sprintf(GET_ERRORS_COUNT_BY_SESSION_LIMIT, filename + "_errors")
		err := r.db.QueryRow(completeQuery).Scan(&count)
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

func (r *LogsDbSqlite) FindLogs(inputs models.FindLogsStruct) ([]models.LogStruct, error) {

	query := fmt.Sprintf(FIND_LOG, inputs.Filename+"_fts")
	inputs.SearchText = fmt.Sprintf("'\"%s\"*';", inputs.SearchText)
	completeQuery := query + inputs.SearchText
	rows, err := r.db.Query(completeQuery)
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

func (r *LogsDbSqlite) GetLogById(id int, filename string) (models.LogStruct, error) {

	var logRow models.LogStruct

	query := fmt.Sprintf(GET_LOG_BY_ID, filename+"_logs")
	err := r.db.QueryRow(query, id).Scan(&logRow.Id, &logRow.SessionId, &logRow.Date, &logRow.ServiceInfo, &logRow.Description)
	if err != nil && err != sql.ErrNoRows {
		tools.LogErr(erx.New(err))
		return logRow, erx.New(err)
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

func (r *LogsDbSqlite) GetLogs() ([]models.LogStruct, error) {

	rows, err := r.db.Query(GET_LOGS)
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

func (r *LogsDbSqlite) GetLogsSessions(filename string) ([]models.LogSessionStruct, error) {

	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	query := fmt.Sprintf(GET_LOGS_SESSIONS, filename+"_sessions")
	rows, err := r.db.Query(query)
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

func (r *LogsDbSqlite) GetLogsServiceInfo(filename string) ([]models.LogStruct, error) {


	if filename == "" {
		return nil, erx.NewError(0, "Error: filename is not specified")
	}

	query := fmt.Sprintf(GET_LOGS_SERVICE_INFO, filename+"_logs")
	rows, err := r.db.Query(query)
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

func (r *LogsDbSqlite) GetLogsFileNames() ([]models.LogFilenameStruct, error) {

	rows, err := r.db.Query(GET_LOGFILE_NAMES)
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

func (r *LogsDbSqlite) GetLogsFilesInfo() ([]models.LogFileStruct, error) {

	rows, err := r.db.Query(GET_LOGFILES)
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

func (r *LogsDbSqlite) GetLogsFileInfo(filename string) (models.LogFileStruct, error) {

	var logRow models.LogFileStruct
	err := r.db.QueryRow(GET_LOGFILE, filename).Scan(&logRow.Id, &logRow.LogfileName, &logRow.FileLength, &logRow.LastSessionDate, &logRow.PreviousDate, &logRow.SessionCount)
	if err != nil {
		return models.LogFileStruct{}, err
	}

	return logRow, nil
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
