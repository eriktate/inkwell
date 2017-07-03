package main

import (
	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/dynamo"
	"github.com/lucsky/cuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	var client inkwell.Client
	log.Println("Connecting to dynamo...")
	client = dynamo.NewClient("local", nil)

	blogService := client.BlogService()
	testBlog := inkwell.Blog{
		ID:    cuid.New(),
		Title: "A simple test blog",
	}

	log.Println("Writing blog...")
	if err := blogService.Write(testBlog); err != nil {
		log.WithError(err).Println("Failed to write blog.")
	}

	log.Println("Retrieving blog...")
	getBlog, err := blogService.Get(testBlog.ID)
	if err != nil {
		log.WithError(err).Println("Failed to get blog.")
	}

	log.WithField("Returned Blog", getBlog).Println("Retrieved")
}
