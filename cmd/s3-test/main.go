package main

import (
	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/s3"
	log "github.com/sirupsen/logrus"
)

func main() {
	var client inkwell.Client
	log.Println("Connecting to s3...")
	client = s3.NewClient("dev", nil)

	blogService := client.BlogService()

	blog, err := blogService.Get("eriktate", "sub-test")
}
