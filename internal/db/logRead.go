package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/kennygrant/sanitize"
	"io"
	"monitorlogs/internal/config"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const unixPrefix = "./"
const windowsPrefix = ".\\"

func newLogStruct() models.LogStruct {
	return models.LogStruct{
		Date:        0,
		ServiceInfo: "-",
		Description: "",
	}
}

func (r *LogsDbSqlite) ReadFolder(folderPath string) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if HasSuffix(path, ".log") {
			files = append(files, path)
			return nil
		}
		return nil
	})
	if err != nil {
		tools.LogErr(erx.New(err))
		return
	}
	for _, file := range files {
		fmt.Println(file)
	}

	for _, fileName := range files {
		err = r.Read(switchPrefix() + fileName)
		if err != nil {
			tools.LogErr(erx.New(err))
		}
	}
}

func (r *LogsDbSqlite) ReadCycle(duration string, folderPath string) {
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		tools.LogErr(erx.New(err))
		parsedDuration = 30 * time.Second
	}
	log := fmt.Sprintf("Reading logs every %v seconds\n", parsedDuration.Seconds())
	tools.LogInfo(log)

	var files []string

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		tools.LogErr(erx.New(err))
		return
	}
	for _, file := range files {
		fmt.Println(file)
	}

	for {
		start := time.Now()
		for _, fileName := range files {
			if strings.Contains(fileName, ".log") {
				err = r.Read(switchPrefix() + fileName)
				if err != nil {
					tools.LogErr(erx.New(err))
				}
			}
		}

		t := time.Now()
		elapsed := t.Sub(start)
		if elapsed <= parsedDuration {
			time.Sleep(parsedDuration - elapsed)
		}
	}
}

func (r *LogsDbSqlite) Read(fullFilename string) error {

	stat, err := os.Stat(fullFilename)
	if err != nil {
		tools.LogErr(erx.New(err))
		return err
	}

	filename := splitFilename(fullFilename)

	isNewFile := true
	var isEmptyStart bool
	prevLengthEnv := os.Getenv(filename + "Length")
	prevLength := 0
	if prevLengthEnv != "" {
		isNewFile = false
		prevLength, err = strconv.Atoi(prevLengthEnv)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}
	}
	prevLength64 := int64(prevLength)
	curLength64 := stat.Size()
	var logString models.LogStruct
	var logs []models.LogStruct
	var sessions []models.LogSessionStruct
	var newSession = false
	var prevDateString = os.Getenv(filename + "PreviousDate")
	var prevDate64 int64
	if prevDateString != "" {
		i, _ := strconv.Atoi(prevDateString)
		prevDate64 = int64(i)
	}
	prevDateUtc := tools.FormatUnixToUTC(prevDate64)
	var seekSessionEnd int
	var sessionCount = os.Getenv(filename + "SessionCount")
	sessionId := 1
	if sessionCount != "" {
		isNewFile = false
		sessionId, err = strconv.Atoi(sessionCount)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}
	}

	session := models.LogSessionStruct{Id: sessionId, Dates: prevDateUtc}
	sessions = append(sessions, session)
	if curLength64 > prevLength64 {

		if isNewFile {
			logInfo, err := r.GetLogsFileInfo(filename)
			if err == sql.ErrNoRows {
				tools.LogInfo("New logfile is: " + filename)
				isEmptyStart = true
				newSession = true
				err = nil
			}
			if err != nil {
				tools.LogErr(erx.New(err))
				return erx.New(err)
			}
			if logInfo.LogfileName != "" {
				var logSlice []models.LogFileStruct
				logSlice = append(logSlice, logInfo)
				err = config.SetLogfilesEnv(logSlice)
				if err != nil {
					tools.LogErr(erx.New(err))
				}
				str := fmt.Sprintf("Logfile %v already exists", filename)
				tools.LogInfo(str)
				return erx.NewError(0, "File already exists")
			}
		}

		tools.LogInfo("Reading file: " + fullFilename)

		file, err := os.Open(fullFilename)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}
		defer file.Close()

		_, err = file.Seek(prevLength64, 0)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}
		var buf strings.Builder
		reader := bufio.NewReader(file)
		for {
			line, prefix, err := reader.ReadLine()
			if err == io.EOF {
				fmt.Printf("End Of File: %s", err)
				break
			}
			if err != nil {
				tools.LogErr(err)
				break
			}
			if len(line) == 0 && isEmptyStart == false {
				seekSessionEnd++
				continue
			}

			if seekSessionEnd >= 4 {
				date := tools.FormatUnixToUTC(prevDate64)
				if len(sessions) > 1 {
					sessions[len(sessions)-1].Dates += " - " + date
				}
				if len(sessions) == 1 {
					if session.Dates == "" {
						session.Dates += date
						sessions = append(sessions, session)
					} else {
						session.Dates += " - " + date
						sessions = append(sessions, session)
					}
				}

				newSession = true
				session.Id++
				seekSessionEnd = 0
			}

			if prefix {
				isEmptyStart = false
				buf.Write(line)
				continue
			}
			if !prefix {
				isEmptyStart = false
				seekSessionEnd = 0
				buf.Write(line)
				s := buf.String()
				s = sanitize.HTML(s)
				lineSlice := strings.Split(s, " ")
				date, err := tools.FormatDateToUnix(lineSlice[0])
				if err != nil {
					date = prevDate64
					var newS []string
					newS = append(newS, fmt.Sprint(prevDate64), s)
					s = appendToString(newS)
				}
				if err == nil {
					prevDate64 = date
				}
				if strings.Contains(s, ": ") && err == nil {
					logString = splitLine(s)
					logString.Date = date

				} else {
					logString = writeDefaultLog(s)
					logString.Date = date
				}
				logString.SessionId = session.Id
				logs = append(logs, logString)
				buf.Reset()
			}
			if newSession {
				newSession = false
				session.Dates = tools.FormatUnixToUTC(prevDate64)
				sessions = append(sessions, session)
			}
		}

		lastSessionDate, _ := tools.FormatDateToUnix(session.Dates)

		logfile := models.LogFileStruct{
			Id:              0,
			LogfileName:     filename,
			FileLength:      int(curLength64),
			LastSessionDate: lastSessionDate,
			PreviousDate:    prevDate64,
			SessionCount:    session.Id,
		}

		if isNewFile {

			err = r.CreateLogDatabase(filename)
			if err != nil {
				tools.LogErr(erx.New(err))
				return err
			}

			tools.LogInfo("Begin writing logs to DB")
			rows, err := r.InsertLogs(logs, filename+"_logs")
			if err != nil {
				tools.LogErr(erx.New(err))
				return err
			}
			tools.LogInfo("Complete writing logs to DB")
			fmt.Println("Rows affected: ", rows)

			for _, v := range sessions {
				err = r.InsertLogSession(v, filename)
				if err != nil {
					tools.LogErr(erx.New(err))
					return err
				}

			}

			err = r.InsertLogFileInfo(logfile)
			if err != nil {
				tools.LogErr(erx.New(err))
				return err
			}

			err = config.SetLogfilesEnv([]models.LogFileStruct{logfile})
			if err != nil {
				tools.LogErr(erx.New(err))
				return err
			}

			return nil
		}

		tools.LogInfo("Begin writing logs to DB")
		rows, err := r.InsertLogs(logs, filename+"_logs")
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}
		tools.LogInfo("Complete writing logs to DB")
		fmt.Println("Rows affected: ", rows)

		for _, sessionStruct := range sessions {
			err = r.InsertLogSession(sessionStruct, filename)
			if err != nil {
				tools.LogErr(erx.New(err))
				return err
			}

		}

		err = r.UpdateLogFileInfo(logfile)
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}

		err = config.SetLogfilesEnv([]models.LogFileStruct{logfile})
		if err != nil {
			tools.LogErr(erx.New(err))
			return err
		}

		return nil
	}

	//tools.LogInfo("Nothing to Read. Logs are not updated")

	return nil
}

func splitLine(s string) models.LogStruct {

	logStruct := newLogStruct()

	logParts := strings.Split(s, ": ")
	headParts := strings.Split(logParts[0], " ")
	if len(headParts) > 4 {
		return writeDefaultLog(s)
	}
	logStruct.ServiceInfo = appendToString(headParts[1:])
	if strings.Contains(logStruct.ServiceInfo, "...") {
		s := strings.Split(logStruct.ServiceInfo, "...")
		logStruct.ServiceInfo = s[len(s)-1]
	}
	logTail := logParts[1:]
	desc := appendToString(logTail)
	logStruct.Description = desc

	return logStruct
}

func writeDefaultLog(s string) models.LogStruct {
	var builder strings.Builder
	var ls = newLogStruct()
	splitParts := strings.Split(s, " ")
	descSlice := splitParts[1:]
	for i := range descSlice {
		builder.WriteString(descSlice[i])
		builder.WriteString(" ")
	}
	ls.Description = builder.String()

	return ls
}

func appendToString(slice []string) string {
	var builder strings.Builder
	for i := range slice {
		builder.WriteString(slice[i])
		builder.WriteString(" ")
	}
	return builder.String()
}

func splitFilename(fullFilename string) (filename string) {

	fullFilenameSlice := strings.Split(fullFilename, "\\")
	filenameWithExt := fullFilenameSlice[len(fullFilenameSlice)-1]
	filenameSlice := strings.Split(filenameWithExt, ".")
	filename = filenameSlice[0]

	if filename == "" {
		fullFilenameSlice = strings.Split(fullFilename, "/")
		filenameWithExt = fullFilenameSlice[len(fullFilenameSlice)-1]
		filenameSlice = strings.Split(filenameWithExt, ".")
		filename = filenameSlice[0]
	}

	return filename

}

func switchPrefix() string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return windowsPrefix
	case "linux":
		return unixPrefix
	default:
		return unixPrefix
	}
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
