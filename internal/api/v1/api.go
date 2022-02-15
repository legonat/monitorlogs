package v1

import (
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/auth"
	"monitorlogs/internal/db"
	"monitorlogs/internal/models"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"net/http"
	"strings"
)

func Login(c *gin.Context) {
	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	login := c.PostForm("login")
	password := c.PostForm("password")

	if login == "" || password == "" {
		tools.LogErr(erx.NewError(0, "Empty inputs"))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Not enough inputs"})
		return
	}
	err := db.Check(login, []byte(password), ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Invalid password"})
		return
	}
	err = db.TryDeleteOldSession("Fingerprint", ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to authorise"})
		return
	}
	err = db.CheckSessionsCount(login, ip)
	if err != nil {
		tools.LogErr(erx.New(err))
	}
	rToken, err := auth.CreateRefreshToken()
	err = db.WriteRefreshToken(login, rToken, ua, "Fingerprint", ip, 60)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Unable to authorise"})
		return
	}
	c.SetCookie("rToken", rToken, 60*60*24*60, "/v1", "localhost:5000", false, true)
	c.SetCookie("login", login, 10000, "/v1", "localhost:5000", false, true)
	c.HTML(http.StatusOK, "successfulLogin.html", gin.H{"at": "aToken", "login": login})
}

func Register(c *gin.Context) {
	ip := c.ClientIP()
	login := c.PostForm("login")
	password := c.PostForm("password")
	if login == "" || password == "" {
		tools.LogErr(erx.NewError(0, "Not enough inputs"))
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Not enough inputs"})
		return
	}
	err := db.Register(login, []byte(password), ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(http.StatusUnprocessableEntity, "error.html", gin.H{"error": "User is already registered"})
		return
	}
	c.HTML(http.StatusOK, "successfulRegistration.html", gin.H{})
}

func Logout(c *gin.Context) {
	ip := c.ClientIP()
	rToken, err := c.Cookie("rToken")
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.html", gin.H{"error": "User is not logged in"})
		return
	}
	c.SetCookie("rToken", "", -1, "/v1", "localhost:5000", false, true)
	c.SetCookie("aToken", "", -1, "/v1", "localhost:5000", false, true)
	c.SetCookie("login", "", -1, "/v1", "localhost:5000", false, true)
	err = db.DeleteSession(rToken, ip)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.html", gin.H{"error": "Unable to delete session"})
		return
	}
	c.HTML(200, "successfulLogout.html", gin.H{})

}

func Block(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	_, err = auth.ParseToken(token)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ip := c.ClientIP()
	var blockInputs models.BlockInputs
	err = c.ShouldBindJSON(&blockInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}

	err = db.Block(blockInputs.Login, ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to block user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Unblock(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	_, err = auth.ParseToken(token)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ip := c.ClientIP()
	user := c.PostForm("user")

	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(200, "error.html", gin.H{"error": "Not enough inputs"})
		return
	}
	err = db.Unblock(user, ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(200, "error.html", gin.H{"error": "Unable to unblock user"})
		return
	}
	c.HTML(200, "error.html", gin.H{"success": true})
}

func Auth(c *gin.Context) {

	token, err := c.Cookie("rToken")
	if err != nil {
		tools.LogWarn("No Cookies with rToken")
		c.HTML(200, "login.html", gin.H{})
		return
	}

	login, _, err := db.CheckRefreshToken(token, "Fingerprint", c.ClientIP())
	if err != nil {
		tools.LogWarn("Refresh Token not found")
		c.HTML(200, "login.html", gin.H{})
		return
	}
	aToken, rToken, err := auth.RefreshSession(login, c.Request.UserAgent(), "Fingerprint", c.ClientIP(), 60)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.HTML(200, "error.html", gin.H{"error": "Unable to Refresh token"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"accessToken": aToken, "refreshToken": rToken})
	c.SetCookie("rToken", rToken, 60*60*24*60, "/v1", "localhost:5000", false, true)
	c.SetCookie("aToken", aToken, 60*15, "/v1", "localhost:5000", false, true)
	c.SetCookie("login", login, 10000, "/v1", "localhost:5000", false, true)
	//c.Redirect(302, "/v1")
	c.HTML(http.StatusOK, "successfulLogin.html", gin.H{"at": aToken})
}

func ExitAll(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	token, err := parseHeader(authHeader)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	login, err := auth.ParseToken(token)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	exitInputs := models.ExitInputs{}
	err = c.ShouldBindJSON(&exitInputs)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inputs"})
		return
	}
	ip := c.ClientIP()
	err = db.DeleteAllSessions(login, exitInputs.Fingerprint, ip)
	if err != nil {
		tools.LogErr(erx.New(err))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unable to delete all sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

func parseHeader(authHeader string) (string, error) {

	if authHeader == "" {
		return "", erx.NewError(611, "Empty Header")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", erx.NewError(612, "Incorrect Header")
	}
	if headerParts[0] != "Bearer" {
		return "", erx.NewError(613, "Incorrect Header type. No Bearer")
	}

	return headerParts[1], nil

}
