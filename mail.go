package main

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

// Binding from JSON
type EmailData struct {
	FromName    string            `form:"fromName" json:"fromName" xml:"fromName"  binding:"required"`
	FromAddress string            `form:"fromAddress" json:"fromAddress" xml:"fromAddress" binding:"required"`
	ToParams    map[string]string `form:"toParams" json:"toParams" xml:"toParams" binding:"required"`
	TemplateId  string            `form:"templateId" json:"templateId" xml:"templateId" binding:"required"`
	Params      map[string]string `form:"params" json:"params" xml:"params" binding:"required"`
}

func SendMail(
	fromName string,
	fromAddress string,
	toParams map[string]string,
	templateId string,
	params map[string]string) (err error, mailResponse string) {

	m := mail.NewV3Mail()

	from := mail.NewEmail(fromName, fromAddress)
	m.SetFrom(from)

	m.SetTemplateID(templateId)

	p := mail.NewPersonalization()

	var tos []*mail.Email
	for toParam := range toParams {
		tos = append(tos, mail.NewEmail(toParam, toParams[toParam]))
	}

	p.AddTos(tos...)

	for param := range params {
		p.SetDynamicTemplateData(param, params[param])
	}

	m.AddPersonalizations(p)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(m)
	// TODO: more advanced logging needed here
	if err != nil {
		log.Println(err)
		return err, "failure"
	} else {
		fmt.Println(response.StatusCode)
		if response.Body == "" {
			fmt.Println(response.Body)
		}
	}
	return nil, "success"
}
