package main

import (
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
	params map[string]string) (err error, mailResponse string, status int) {

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

	log.Printf("Sending mail request via sendgird: %v\n", m)
	response, err := client.Send(m)
	log.Printf("Sending mail via sendgrid completed!\n Err: %v,\nResponse: %v\n", err, response)
	if err != nil {
		log.Println(err)
		return err, response.Body, response.StatusCode
	}
	return nil, response.Body, response.StatusCode
}
