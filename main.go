package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/autlamps/delay-backend-api/api"
)

var dburl string
var rdurl string
var key string
var mgKey string
var domain string
var confirmDomain string

func init() {
	flag.StringVar(&dburl, "DATABASE_URL", "", "database url")
	flag.StringVar(&rdurl, "REDIS_URL", "", "redis url")
	flag.StringVar(&key, "KEY", "", "key used to sign jwt")
	flag.StringVar(&mgKey, "MG_KEY", "", "mailgun api key")
	flag.StringVar(&domain, "MG_DOMAIN", "", "mailgun domain")
	flag.StringVar(&confirmDomain, "MG_CONFIRM", "", "delay confirm domain")
	flag.Parse()

	if dburl == "" {
		dburl = os.Getenv("DATABASE_URL")
	}

	if rdurl == "" {
		rdurl = os.Getenv("REDIS_URL")
	}

	if key == "" {
		key = os.Getenv("KEY")
	}

	if mgKey == "" {
		mgKey = os.Getenv("MG_KEY")
	}

	if domain == "" {
		domain = os.Getenv("MG_DOMAIN")
	}

	if confirmDomain == "" {
		confirmDomain = os.Getenv("MG_CONFIRM")
	}
}

func main() {
	c := api.Conf{
		RDURL:         rdurl,
		DBURL:         dburl,
		Key:           key,
		MGKey:         mgKey,
		Domain:        domain,
		ConfirmDomain: confirmDomain,
	}

	r, err := api.Create(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":5000", r))
}
