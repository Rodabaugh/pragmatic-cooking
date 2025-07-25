package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

func (cfg *apiConfig) sendMGEmail(name, rcpt, subject, msg string) {
	var domain = "noreply.pragmatic.cooking"

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(cfg.mgAPIKey)

	sender := "noreply@noreply.pragmatic.cooking"

	// The message object allows you to add attachments and Bcc recipients
	message := mailgun.NewMessage(domain, sender, subject, msg, rcpt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10-second timeout
	resp, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", resp.ID, resp.Message)
}
