package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mailgun/mailgun-go/v4"
)

func main() {
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunAPIKey := os.Getenv("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun(mailgunDomain, mailgunAPIKey)

	sender := fmt.Sprintf("mailgun@%s", os.Getenv("MAILGUN_DOMAIN"))
	subject := "Cloud Computing Homework - Email"
	body := "Hello from Mailgun Go!"
	recipient := os.Getenv("MAILGUN_RECIPIENT")

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
