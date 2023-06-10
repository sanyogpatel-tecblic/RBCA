package email

import (
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailAlert(email string, apiinfo string, Message string) {
	sendgridAPIKey := "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"

	from := mail.NewEmail("admin", "ad2491min@gmail.com")
	subject := "API Alert"
	to := mail.NewEmail("sanyog", email)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")
	content := mail.NewContent("text/plain", apiinfo+Message+" at time: "+timestamp)
	message := mail.NewV3MailInit(from, subject, to, content)

	// Create a new SendGrid client
	client := sendgrid.NewSendClient(sendgridAPIKey)

	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send email alert:", err)
		return
	}
	// Check the response status code
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		log.Println("Email alert sent successfully")
	} else {
		log.Println("Failed to send email alert. Status code:", response.StatusCode)
	}
}

func SendEmailAlert2(fromemail string, email string, apiinfo string, Details string) {
	sendgridAPIKey := "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"

	timestamp := time.Now().Format("2006-01-02 15:04:05.000000")

	from := mail.NewEmail("admin", fromemail)
	subject := "API Alert"
	to := mail.NewEmail("admin", email)
	content := mail.NewContent("text/plain", apiinfo+" with username "+Details+" at time: "+timestamp)
	message := mail.NewV3MailInit(from, subject, to, content)

	// Create a new SendGrid client
	client := sendgrid.NewSendClient(sendgridAPIKey)

	// Send the email using the SendGrid client
	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send email alert:", err)
		return
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		log.Println("Email alert sent successfully")
	} else {
		log.Println("Failed to send email alert. Status code:", response.StatusCode)
	}
}
