package s3

import (
	"os"

	"github.com/aws-sdk-go/aws/credentials"
	"github.com/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/eriktate/inkwell"
)

type Client struct {
	svc         s3iface.S3API
	blogService BlogService
}

func NewClient(env string, svc s3iface.S3API) *Client {
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

func (c *Client) BlogService() inkwell.BlogService {
	return c.blogService
}
