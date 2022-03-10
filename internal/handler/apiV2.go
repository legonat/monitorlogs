package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/auth"
	"monitorlogs/internal/config"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) Login(c *gin.Context) {
	var loginInputs models.LoginInputs
	err := c.ShouldBindJSON(&loginInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	fingerprint := c.GetHeader("Fingerprint")
	if fingerprint == "" {
		tools.LogErr(erx.NewError(0, "No fingerprint data"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	loginInputs.Ip = c.ClientIP()
	err = h.repository.Check(loginInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	err = h.repository.TryDeleteOldSession(fingerprint, loginInputs.Ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to authorise"})
		return
	}
	err = h.repository.CheckSessionsCount(loginInputs.Login, loginInputs.Ip)
	if err != nil {
		tools.LogErr(erx.New(err))
	}

	refreshSession :=  models.RefreshSession{Login: loginInputs.Login, Ip: loginInputs.Ip, Fingerprint: fingerprint}
	refreshSession.Token, err = auth.CreateRefreshToken()
	refreshSession.Ua = c.Request.UserAgent()
	domain, daysUntilExpire, err := getConfig()
	if !loginInputs.RememberMe {
		daysUntilExpire = 1
	}
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to Refresh token"})
		return
	}
	err = h.repository.WriteRefreshToken(refreshSession, daysUntilExpire)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to authorise"})
		return
	}

	c.SetCookie("rToken", refreshSession.Token, daysUntilExpire*60*60*24, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "refreshToken": refreshSession.Token})
}

func (h *Handler) Register(c *gin.Context) {
	var registerInputs models.RegisterInputs
	err := c.ShouldBindJSON(&registerInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	registerInputs.Ip = c.ClientIP()
	err = h.repository.Register(registerInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User is already registered"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Block(c *gin.Context) {

	var blockInputs models.BlockInputs
	err := c.ShouldBindJSON(&blockInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	blockInputs.Ip = c.ClientIP()

	err = h.repository.Block(blockInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to block user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Unblock(c *gin.Context) {

	var unblockInputs models.BlockInputs
	err := c.ShouldBindJSON(&unblockInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	unblockInputs.Ip= c.ClientIP()
	err = h.repository.Unblock(unblockInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to unblock user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Auth(c *gin.Context) {

	var inputs models.RefreshSession
	inputs.Fingerprint = c.GetHeader("Fingerprint")

	domain, _, err := getConfig()
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	inputs.Token, err = c.Cookie("rToken")
	if err != nil {
		tools.LogWarn("No Cookies")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	inputs.Ip = c.ClientIP()
	inputs.Ua = c.Request.UserAgent()
	var daysUntilExpire int
	inputs.Login, daysUntilExpire, err = h.repository.CheckRefreshToken(inputs)
	if err != nil {
		tools.LogErr(err)
		if err == sql.ErrNoRows {
			c.SetCookie("rToken", "", -1, "/", domain, false, true)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	aToken, rToken, err := auth.RefreshSession(inputs, daysUntilExpire)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	inputs.Token = rToken

	err = h.repository.WriteRefreshToken(inputs, daysUntilExpire)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}

	c.SetCookie("rToken", rToken, daysUntilExpire*60*60*24, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "accessToken": aToken, "refreshToken": rToken})
}

func (h *Handler) Logout(c *gin.Context) {

	refreshToken, err := c.Cookie("rToken")
	if err != nil {
		tools.LogWarn("No Cookies")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ip := c.ClientIP()
	err = h.repository.DeleteSession(refreshToken, ip)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to delete session"})
		return
	}
	domain, _, err := getConfig()
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to Refresh token"})
		return
	}
	c.SetCookie("rToken", "", -1, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ExitAll(c *gin.Context) {
	var inputs models.RefreshSession
	var err error
	authHeader := c.GetHeader("Authorization")
	inputs.Token, err = parseHeader(authHeader)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	inputs.Login, err = auth.ParseToken(inputs.Token)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	inputs.Ip = c.ClientIP()
	inputs.Fingerprint = c.GetHeader("Fingerprint")

	err = h.repository.DeleteAllSessions(inputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to delete all sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

func (h *Handler) GetLogs(c *gin.Context) {

	logs, err := h.repository.GetLogs()
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *Handler) GetLogsBySession(c *gin.Context) {

	reqInputs := models.GetLogsBySessionStruct{}
	err := c.BindJSON(&reqInputs)
	logs, err := h.repository.GetLogsBySession(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, err := h.repository.GetErrorsBySession(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "errors": errors})
}

func (h *Handler) GetLogsByDate(c *gin.Context) {

	reqInputs := models.GetLogsByDateStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools.LogErr(erx.New(err))
		return
	}
	dif := tools.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)
	logs, err := h.repository.GetLogsByDate(startDate, endDate, reqInputs.Filename)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, err := h.repository.GetErrorsByDate(startDate, endDate, reqInputs.Filename)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "errors": errors})
}

func (h *Handler) GetLogsBySessionWithLimit(c *gin.Context) {

	logSlices := [3][]models.LogStruct{}

	reqInputs := models.GetLogsBySessionWithLimitStruct{}
	err := c.BindJSON(&reqInputs)
	logs, logsCount, err := h.repository.GetLogsBySessionWithLimit(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := h.repository.GetErrorsBySessionWithLimit(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "errorsSlice": errors[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func (h *Handler) GetLogsByDateWithLimit(c *gin.Context) {
	logSlices := [3][]models.LogStruct{}
	reqInputs := models.GetLogsByDateWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools.LogErr(erx.New(err))
		return
	}
	dif := tools.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)
	logs, logsCount, err := h.repository.GetLogsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := h.repository.GetErrorsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}
	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "errorsSlice": errors[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func (h *Handler) GetErrorsBySessionWithLimit(c *gin.Context) {
	errorSlices := [3][]models.LogStruct{}
	reqInputs := models.GetLogsBySessionWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	logs, logsCount, err := h.repository.GetLogsBySessionWithLimit(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := h.repository.GetErrorsBySessionWithLimit(reqInputs)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get errors"})
		return
	}
	for i, slice := range errors {
		errorSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backErrors": errorSlices[0], "currentErrors": errorSlices[1], "forwardErrors": errorSlices[2], "logsSlice": logs[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func (h *Handler) GetErrorsByDateWithLimit(c *gin.Context) {

	errorSlices := [3][]models.LogStruct{}
	reqInputs := models.GetLogsByDateWithLimitStruct{}
	err := c.BindJSON(&reqInputs)

	start, err := time.Parse(time.RFC3339, reqInputs.StartDate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to parse date"})
		tools.LogErr(erx.New(err))
		return
	}
	dif := tools.TimeFix(start.Hour())
	startDate := start.Add(dif)
	end, err := time.Parse(time.RFC3339, reqInputs.EndDate)
	if err != nil {
		tools.LogErr(erx.New(err))
		end = time.Unix(0, 0)
	}
	dif = tools.TimeFix(end.Hour())
	endDate := end.Add(dif)
	endDate = endDate.AddDate(0, 0, 1)

	logs, logsCount, err := h.repository.GetLogsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get logs"})
		return
	}

	errors, errorsCount, err := h.repository.GetErrorsByDateWithLimit(startDate, endDate, reqInputs.Filename, reqInputs.Limit, reqInputs.Offset)
	if err != nil {
		tools.LogErr(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get errors"})
		return
	}
	for i, slice := range errors {
		errorSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backErrors": errorSlices[0], "currentErrors": errorSlices[1], "forwardErrors": errorSlices[2], "logsSlice": logs[1], "logsCount": logsCount, "errorsCount": errorsCount})
}

func (h *Handler) FindLogs(c *gin.Context) {

	reqInputs := models.FindLogsStruct{}
	err := c.BindJSON(&reqInputs)

	logs, err := h.repository.FindLogs(reqInputs)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find logs"})
		tools.LogErr(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

func (h *Handler) FindLogsWithLimit(c *gin.Context) {

	logSlices := [3][]models.LogStruct{}
	reqInputs := models.FindLogsStructWithLimit{}
	err := c.BindJSON(&reqInputs)

	logs, logsCount, err := h.repository.FindLogsWithLimit(reqInputs)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find logs"})
		tools.LogErr(err)
		return
	}

	for i, slice := range logs {
		logSlices[i] = slice
	}

	c.JSON(http.StatusOK, gin.H{"backLogs": logSlices[0], "currentLogs": logSlices[1], "forwardLogs": logSlices[2], "logsCount": logsCount})
}

func (h *Handler) GetLogById(c *gin.Context) {

	reqInputs := models.LogFilenameStruct{}
	err := c.ShouldBindJSON(&reqInputs)

	log, err := h.repository.GetLogById(reqInputs.Id, reqInputs.LogfileName)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to find log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": log})

}

func (h *Handler) GetLogsSessions(c *gin.Context) {

	reqInputs := models.LogFilenameStruct{}
	err := c.BindJSON(&reqInputs)

	sessions, err := h.repository.GetLogsSessions(reqInputs.LogfileName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get sessions list"})
		tools.LogErr(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *Handler) GetLogsServiceInfo(c *gin.Context) {

	reqInputs := models.LogFilenameStruct{}
	err := c.BindJSON(&reqInputs)

	services, err := h.repository.GetLogsServiceInfo(reqInputs.LogfileName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to get sessions list"})
		tools.LogErr(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"services": services})
}

func (h *Handler) GetLogsFilenames(c *gin.Context) {

	filenames, err := h.repository.GetLogsFileNames()
	if err != nil {
		tools.LogErr(err)
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
		tools.LogErr(erx.New(err))
	}
	if domain == "" || daysUntilExpire == 0 {
		cfg := config.GetInstance()
		if err != nil {
			tools.LogErr(erx.New(err))
			return "", 0, err
		}
		domain = (*cfg).Server.Domain
		daysUntilExpire = (*cfg).Server.MaxAge
	}
	if domain == "" || daysUntilExpire == 0 {
		cfg, err := config.GetConfig()
		if err != nil {
			tools.LogErr(erx.New(err))
			return "", 0, err
		}
		domain = cfg.Server.Domain
		daysUntilExpire = cfg.Server.MaxAge
	}
	return domain, daysUntilExpire, err
}
