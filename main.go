package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func setupRouter() *gin.Engine {
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Ping test
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I am busy!",
		})
	})

	router.POST("sendMail", func(c *gin.Context) {
		var emailData EmailData
		if err := c.ShouldBindJSON(&emailData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err, response := SendMail(
			emailData.FromName,
			emailData.FromAddress,
			emailData.ToParams,
			emailData.TemplateId,
			emailData.Params)
		if err == nil {
			c.JSON(200, gin.H{
				"message": response,
			})
		} else {
			c.JSON(400, gin.H{
				"message": response,
			})
		}

	})

	router.POST("sendSms", func(c *gin.Context) {
		var smsData SMSData
		if err := c.ShouldBindJSON(&smsData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err, response := SendSMS(
			smsData.FromNumber,
			smsData.ToNumber,
			smsData.Message)
		if err == nil {
			c.JSON(200, gin.H{
				"message": response,
			})
		} else {
			c.JSON(400, gin.H{
				"message": response,
			})
		}

	})

	return router
}

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})

	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	router := setupRouter()
	_ = router.Run(":7070")
}
