package email

import (
	"flag"
	"os"
	"testing"
)

var key string
var domain string
var confirmDomain string

func init() {
	flag.StringVar(&key, "MG_KEY", "", "mailgun api key")
	flag.StringVar(&domain, "MG_DOMAIN", "", "mailgun domain")
	flag.StringVar(&confirmDomain, "MG_CONFIRM", "", "delay confirm domain")
	flag.Parse()

	if key == "" {
		key = os.Getenv("MG_KEY")
	}

	if domain == "" {
		domain = os.Getenv("MG_DOMAIN")
	}

	if confirmDomain == "" {
		confirmDomain = os.Getenv("MG_CONFIRM")
	}
}

func TestMailgun_SendConfirmation(t *testing.T) {
	mg := InitMailgunService(key, domain, confirmDomain)

	err := mg.SendConfirmation("hayden@example.co.nz", "Hayden Woodhead", "1234")

	if err != nil {
		t.Fatalf("Email failed: %v", err)
	}
}
