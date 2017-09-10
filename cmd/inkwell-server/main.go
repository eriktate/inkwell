package main

import (
	"github.com/eriktate/inkwell/dynamo"
	"github.com/eriktate/inkwell/http"
	"github.com/eriktate/inkwell/s3"
	"github.com/eriktate/inkwell/sdyn"
	"github.com/pressly/chi"
	log "github.com/sirupsen/logrus"
)

func main() {
	server := http.NewServer("localhost", 8080)
	s3Client := s3.NewClient("dev", nil)
	dynClient := dynamo.NewClient("dev", nil)

	sdynClient := sdyn.NewClient(s3Client, dynClient)

	blogHandler := http.NewBlogHandler(sdynClient.BlogService())

	router := chi.NewRouter()
	router.Get("/health", http.HealthCheck)
	router.Get("/author/{authorID}/blog/{blogID}", blogHandler.Get)
	router.Get("/page/author/{authorID}/blog/{blogID}", blogHandler.GetPage)
	router.Post("/blog", blogHandler.Write)
	router.Post("/author/{authorID}/blog/{blogID}", blogHandler.Revise)
	router.Put("/author/{authorID}/blog/{blogID}", blogHandler.Publish)
	router.Delete("/author/{authorID}/blog/{blogID}", blogHandler.Redact)

	if err := server.Start(router); err != nil {
		log.WithError(err).Println("Server failed")
	}
}
