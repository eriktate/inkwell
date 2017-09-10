package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/eriktate/inkwell"
	log "github.com/sirupsen/logrus"
)

// Client provides an S3 implemention of the inkwell.Client interface.
type Client struct {
	svc         s3iface.S3API
	blogService BlogService
}

// NewClient returns a new S3 inkwell client.
func NewClient(env string, svc s3iface.S3API) *Client {
	log.Println("Dialing s3...")
	if svc == nil {
		if env == "local" {
			localS3 := os.Getenv("LOCAL_S3")
			log.WithField("Local s3", localS3).Println("Connecting to s3...")
			svc = s3.New(session.New(&aws.Config{
				Endpoint:    aws.String(localS3),
				Credentials: credentials.NewStaticCredentials("local", "test", "stuff"),
				Region:      aws.String("us-east-1"),
			}))
		} else {
			svc = s3.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
		}
	}

	return &Client{
		svc:         svc,
		blogService: NewBlogService(svc, os.Getenv("INKWELL_BLOGS_BUCKET")),
	}
}

// BlogService fulfills the inwekll.Client interface and provides an S3
// implementation of an inkwell.BlogService
func (c *Client) BlogService() inkwell.BlogService {
	return c.blogService
}
