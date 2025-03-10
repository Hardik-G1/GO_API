package mailer

import (
	"context"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var yourDomain string = "mg.domain.com" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var privateAPIKey string = "{API_KEY}" //

func SendAccountConfirmation(code string, email string) {
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun("mg.formulas.tech", privateAPIKey)

	sender := "admin@formulas.com"
	subject := "Do not Reply - Account Confirmation"
	body := "This is about the account confirmation on formulas the code for activation is " + code + " . ONLY VALID FOR 10 MIN."
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("ID: %s Resp: %s\n", id, resp)
}
