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
	testBlog := inkwell.Blog{
		ID:        "first-post",
		AuthorID:  "thediylife",
		Title:     "First Post",
		Published: true,
		Content: `
			<html>
			<head>
			<title>thediylife - First Post</title>
			<body>
				<h1>Welcome to thediylife blog!</h1>
				<p>Insert awesome project here!</p>
			</body>
			</head>
			</html>
		`,
	}

	log.Println("Writing blog...")
	if err := blogService.Write(testBlog); err != nil {
		log.WithError(err).Println("Failed to write blog.")
	}
}
