package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))
	r.Use(adapter.Wrap(csrfMiddleware))

	r.LoadHTMLGlob("templates/*/*")

	r.GET("/", pageIndex)
	r.POST("/", pageIndexPOST)

	r.GET("/about", pageAbout)

	r.POST("/api/test", apiTest)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":9002")
}

func pageIndex(c *gin.Context) {

	c.Writer.Header().Set("X-CSRF-Token", csrf.Token(c.Request))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"csrf": csrf.TemplateField(c.Request),
	})
}

func pageAbout(c *gin.Context) {

	c.HTML(http.StatusOK, "about.html", nil)
}

func pageIndexPOST(c *gin.Context) {
	username := c.PostForm("username")

	if strings.TrimSpace(username) == "" {
		c.HTML(http.StatusOK, "signin.html", gin.H{
			"title":   "Sign In Failed",
			"csrf":    csrf.TemplateField(c.Request),
			"message": "Email & Password Required",
		})
		return
	}

	c.Writer.Header().Set("X-CSRF-Token", csrf.Token(c.Request))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"csrf":    csrf.TemplateField(c.Request),
		"message": fmt.Sprintf("form post: %s", username),
	})
}

type reqTest struct {
	Text string `json:"text"`
}

func apiTest(c *gin.Context) {
	var req reqTest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("api post: %s", req.Text)})
}
