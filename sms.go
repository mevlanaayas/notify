package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Binding from JSON
type SMSData struct {
	FromNumber string `form:"fromNumber" json:"fromNumber" xml:"fromNumber" binding:"required"`
	ToNumber   string `form:"toNumber" json:"toNumber" xml:"toNumber" binding:"required"`
	Message    string `form:"message" json:"message" xml:"message" binding:"required"`
}

func SendSMS(
	fromNumber string,
	toNumber string,
	message string) (err error, smsResponse string) {

	msgData := url.Values{}
	msgData.Set("To", toNumber)
	msgData.Set("From", fromNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}

	accountSid := os.Getenv("LIVE_ACCOUNT_SID")
	authToken := os.Getenv("LIVE_AUTH_TOKEN")

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return err, "failure"
	} else {
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			var data map[string]interface{}
			decoder := json.NewDecoder(response.Body)
			err := decoder.Decode(&data)
			if err == nil {
				fmt.Println(data["sid"])
			}
		} else {
			fmt.Println(response.Status)
		}
	}
	return nil, "success"
}
