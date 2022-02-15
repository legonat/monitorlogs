package v2

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/auth"
	"monitorlogs/internal/config"
	"monitorlogs/internal/db"
	models2 "monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	tools2 "monitorlogs/pkg/tools"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Login(c *gin.Context) {
	var checkInputs models2.LoginInputs
	err := c.ShouldBindJSON(&checkInputs)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	fingerprint := c.GetHeader("Fingerprint")
	if fingerprint == "" {
		tools2.LogErr(erx.NewError(0, "No fingerprint data"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	err = db.Check(checkInputs.Login, []byte(checkInputs.Password), ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	err = db.TryDeleteOldSession(fingerprint, ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to authorise"})
		return
	}
	err = db.CheckSessionsCount(checkInputs.Login, ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
	}

	rToken, err := auth.CreateRefreshToken()
	domain, daysUntilExpire, err := getConfig()
	if !checkInputs.RememberMe {
		daysUntilExpire = 1
	}
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to Refresh token"})
		return
	}
	err = db.WriteRefreshToken(checkInputs.Login, rToken, ua, fingerprint, ip, daysUntilExpire)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to authorise"})
		return
	}

	c.SetCookie("rToken", rToken, daysUntilExpire*60*60*24, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "refreshToken": rToken})
}

func Register(c *gin.Context) {
	var checkInputs models2.RegisterInputs
	err := c.ShouldBindJSON(&checkInputs)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	ip := c.ClientIP()
	err = db.Register(checkInputs.Login, []byte(checkInputs.Password), ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User is already registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Block(c *gin.Context) {

	ip := c.ClientIP()
	var blockInputs models2.BlockInputs
	err := c.ShouldBindJSON(&blockInputs)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}

	err = db.Block(blockInputs.Login, ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to block user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Unblock(c *gin.Context) {

	ip := c.ClientIP()
	var unblockInputs models2.BlockInputs
	err := c.ShouldBindJSON(&unblockInputs)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	err = db.Unblock(unblockInputs.Login, ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to unblock user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Auth(c *gin.Context) {

	fingerprint := c.GetHeader("Fingerprint")

	refreshToken, err := c.Cookie("rToken")
	if err != nil {
		tools2.LogWarn("No Cookies")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	domain, _, err := getConfig()
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	login, daysUntilExpire, err := db.CheckRefreshToken(refreshToken, fingerprint, ip)
	if err != nil {
		tools2.LogErr(err)
		if err == sql.ErrNoRows {
			c.SetCookie("rToken", "", -1, "/", domain, false, true)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	aToken, rToken, err := auth.RefreshSession(login, ua, fingerprint, ip, daysUntilExpire)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	c.SetCookie("rToken", rToken, daysUntilExpire*60*60*24, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "accessToken": aToken, "refreshToken": rToken})
}

func Logout(c *gin.Context) {

	refreshToken, err := c.Cookie("rToken")
	if err != nil {
		tools2.LogWarn("No Cookies")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ip := c.ClientIP()
	err = db.DeleteSession(refreshToken, ip)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to delete session"})
		return
	}
	domain, _, err := getConfig()
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}
	c.SetCookie("rToken", "", -1, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ExitAll(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	login, err := auth.ParseToken(token)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//
	//exitInputs := models.ExitInputs{}
	//err = c.ShouldBindJSON(&exitInputs)
	//if err != nil {
	//	tools.LogErr(erx.New(err))
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
	//	return
	//}
	ip := c.ClientIP()
	fingerprint := c.GetHeader("Fingerprint")

	err = db.DeleteAllSessions(login, fingerprint, ip)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to delete all sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

func GetLogs(c *gin.Context) {

	logs, err := db.GetLogs()
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func GetLogsBySession(c *gin.Context) {

	reqInputs := models2.GetLogsBySessionStruct{}
	err := c.BindJSON(&reqInputs)
	logs, err := db.GetLogsBySession(reqInputs.SessionId, reqInputs.Filename)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, err := db.GetErrorsBySession(reqInputs.SessionId, reqInputs.Filename)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "errors": errors})
}

func GetLogsByDate(c *gin.Context) {

	reqInputs := models2.GetLogsByDateStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools2.LogErr(erx.New(err))
		return
	}
	dif := tools2.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools2.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools2.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)
	logs, err := db.GetLogsByDate(startDate, endDate, reqInputs.Filename)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, err := db.GetErrorsByDate(startDate, endDate, reqInputs.Filename)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "errors": errors})
}

func GetLogsBySessionWithLimit(c *gin.Context) {

	logSlices := [3][]models2.LogStruct{}

	reqInputs := models2.GetLogsBySessionWithLimitStruct{}
	err := c.BindJSON(&reqInputs)
	logs, logsCount, err := db.GetLogsBySessionWithLimit(reqInputs.SessionId, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := db.GetErrorsBySessionWithLimit(reqInputs.SessionId, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "errorsSlice": errors[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func GetLogsByDateWithLimit(c *gin.Context) {
	logSlices := [3][]models2.LogStruct{}
	reqInputs := models2.GetLogsByDateWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools2.LogErr(erx.New(err))
		return
	}
	dif := tools2.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools2.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools2.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)
	logs, logsCount, err := db.GetLogsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := db.GetErrorsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}
	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "errorsSlice": errors[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func GetErrorsBySessionWithLimit(c *gin.Context) {
	errorSlices := [3][]models2.LogStruct{}
	reqInputs := models2.GetLogsBySessionWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	logs, logsCount, err := db.GetLogsBySessionWithLimit(reqInputs.SessionId, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := db.GetErrorsBySessionWithLimit(reqInputs.SessionId, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get errors"})
		return
	}
	for i, slice := range errors {
		errorSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backErrors": errorSlices[0], "currentErrors": errorSlices[1], "forwardErrors": errorSlices[2], "logsSlice": logs[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func GetErrorsByDateWithLimit(c *gin.Context) {

	errorSlices := [3][]models2.LogStruct{}
	reqInputs := models2.GetLogsByDateWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools2.LogErr(erx.New(err))
		return
	}
	dif := tools2.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools2.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools2.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)

	logs, logsCount, err := db.GetLogsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := db.GetErrorsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get errors"})
		return
	}
	for i, slice := range errors {
		errorSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backErrors": errorSlices[0], "currentErrors": errorSlices[1], "forwardErrors": errorSlices[2], "logsSlice": logs[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func FindLogs(c *gin.Context) {

	reqInputs := models2.FindLogsStruct{}
	err := c.BindJSON(&reqInputs)

	logs, err := db.FindLogs(reqInputs.SearchText, reqInputs.Filename)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find logs"})
		tools2.LogErr(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func FindLogsWithLimit(c *gin.Context) {

	logSlices := [3][]models2.LogStruct{}
	reqInputs := models2.FindLogsStructWithLimit{}
	err := c.BindJSON(&reqInputs)

	logs, logsCount, err := db.FindLogsWithLimit(reqInputs.SearchText, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find logs"})
		tools2.LogErr(err)
		return
	}

	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "logsCount": logsCount})
}

func GetLogById(c *gin.Context) {

	reqInputs := models2.LogFilenameStruct{}
	err := c.ShouldBindJSON(&reqInputs)

	log, err := db.GetLogById(reqInputs.Id, reqInputs.LogfileName)
	if err != nil {
		tools2.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": log})

}

func GetLogsSessions(c *gin.Context) {

	reqInputs := models2.LogFilenameStruct{}
	err := c.BindJSON(&reqInputs)

	sessions, err := db.GetLogsSessions(reqInputs.LogfileName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get sessions list"})
		tools2.LogErr(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func GetLogsServiceInfo(c *gin.Context) {

	reqInputs := models2.LogFilenameStruct{}
	err := c.BindJSON(&reqInputs)

	services, err := db.GetLogsServiceInfo(reqInputs.LogfileName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get sessions list"})
		tools2.LogErr(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"services": services})
}

func GetLogsFilenames(c *gin.Context) {

	filenames, err := db.GetLogsFileNames()
	if err != nil {
		tools2.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get sessions list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": filenames})
}

func parseHeader(header string) (string, error) {

	if header == "" {
		return "", erx.NewError(611, "Empty Header")
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", erx.NewError(612, "Incorrect Header")
	}
	if headerParts[0] != "Bearer" {
		return "", erx.NewError(613, "Incorrect Header type. No Bearer")
	}

	return headerParts[1], nil

}

func getConfig() (string, int, error) {
	domain := ""
	domain = os.Getenv("DOMAIN")
	//daysUntilExpire := 0
	daysUntilExpire, err := strconv.Atoi(os.Getenv("COOKIE_MAX_AGE"))
	if err != nil {
		tools2.LogErr(erx.New(err))
	}
	if domain == "" || daysUntilExpire == 0 {
		cfg := config.GetInstance()
		if err != nil {
			tools2.LogErr(erx.New(err))
			return "", 0, err
		}
		domain = (*cfg).Server.Domain
		daysUntilExpire = (*cfg).Server.MaxAge
	}
	if domain == "" || daysUntilExpire == 0 {
		cfg, err := config.GetConfig()
		if err != nil {
			tools2.LogErr(erx.New(err))
			return "", 0, err
		}
		domain = cfg.Server.Domain
		daysUntilExpire = cfg.Server.MaxAge
	}
	return domain, daysUntilExpire, err
}
