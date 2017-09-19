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

func init() {
	flag.StringVar(&dburl, "DB_URL", "", "database url")
	flag.StringVar(&rdurl, "RD_URL", "", "redis url")
	flag.Parse()

	if dburl == "" {
		dburl = os.Getenv("DB_URL")
	}

	if rdurl == "" {
		rdurl = os.Getenv("RD_URL")
	}
}

func main() {
	c := api.Conf{"", dburl, "hello"}

	r, err := api.Create(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
