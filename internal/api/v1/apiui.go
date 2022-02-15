package v1

import "github.com/gin-gonic/gin"

func ShowMain(c *gin.Context) {

	_, err := c.Cookie("rToken")
	if err != nil {
		c.Redirect(302, "/v1/login")
		return
	}

	login, err := c.Cookie("login")
	if err != nil {
		//c.HTML(200, "login.html", gin.H{})
		//c.Request.URL.Path = "/v1/login"
		c.Redirect(302, "/v1/login")
		return
	}
	c.HTML(200, "successfulLogin.html", gin.H{"login": login})
}

func ShowRegistrationPage(c *gin.Context) {

	c.HTML(200, "registration.html", gin.H{})

}

func ShowLoginPage(c *gin.Context) {

	c.HTML(200, "login.html", gin.H{})

}
