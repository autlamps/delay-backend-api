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

func init() {
	flag.StringVar(&dburl, "DATABASE_URL", "", "database url")
	flag.StringVar(&rdurl, "REDIS_URL", "", "redis url")
	flag.StringVar(&key, "KEY", "", "key used to sign jwt")
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
}

func main() {
	c := api.Conf{rdurl, dburl, key}

	r, err := api.Create(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":5000", r))
}
