package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	message := url.Values{}
	message.Set("To", os.Getenv("TARGET_NUMBER"))
	message.Set("From", "12058430161")
	message.Set("Body", "Your Mirror of Kalandra is ready for pickup")

	reader := *strings.NewReader(message.Encode())

	client := &http.Client{}

	req, _ := http.NewRequest("POST", endpoint, &reader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
}
