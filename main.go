package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync"
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
		defer recovery(c, nil)
	})

	router.POST("sendMail", func(c *gin.Context) {
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)

		go func() {
			defer recovery(c, &waitGroup)
			var emailData EmailData
			if err := c.ShouldBindJSON(&emailData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				waitGroup.Done()
				return
			}
			log.Printf("[Send Email] Mail request to user %s, with params: %s\n", emailData.ToParams, emailData.Params)
			err, response, status := SendMail(
				emailData.FromName,
				emailData.FromAddress,
				emailData.ToParams,
				emailData.TemplateId,
				emailData.Params)
			if err == nil {
				c.JSON(status, gin.H{
					"message": response,
				})
				waitGroup.Done()
				return
			} else {
				c.JSON(status, gin.H{
					"message": response,
				})
				waitGroup.Done()
				return
			}
		}()
		waitGroup.Wait()
	})

	router.POST("sendSms", func(c *gin.Context) {
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)
		go func() {
			defer recovery(c, nil)
			var smsData SMSData
			if err := c.ShouldBindJSON(&smsData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				waitGroup.Done()
				return
			}
			log.Printf("[Send Sms] Sms request to user %s, with params: %s\n", smsData.ToNumber, smsData.Message)
			err, response, status := SendSMS(
				smsData.FromNumber,
				smsData.ToNumber,
				smsData.Message)
			if err == nil {
				c.JSON(status, gin.H{
					"message": response,
				})
				waitGroup.Done()
				return
			} else {
				c.JSON(status, gin.H{
					"message": response,
				})
				waitGroup.Done()
				return
			}
		}()
		waitGroup.Wait()
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

func recovery(context *gin.Context, waitGroup *sync.WaitGroup) {
	if r := recover(); r != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Result": "[unexpected Notify error]. Recovered!"})
		log.Println("Recovered ", r)
		if waitGroup != nil {
			waitGroup.Done()
		}
	}
}
