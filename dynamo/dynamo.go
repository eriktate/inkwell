package dynamo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/eriktate/inkwell"
	log "github.com/sirupsen/logrus"
)

// A Client wraps the connection to dynamo and provides service functions
// for working with dynamo data from Inkwell.
type Client struct {
	db          dynamodbiface.DynamoDBAPI
	blogService BlogService
}

// NewClient returns a new dynamo Client configured with the given env
// and optional existing dynamo db.
func NewClient(env string, db dynamodbiface.DynamoDBAPI) *Client {
	log.Println("Dialing dynamo...")
	if db == nil {
		if env == "local" {
			localDynamo := os.Getenv("LOCAL_DYNAMO")
			log.WithField("Local Dynamo", localDynamo).Println("Connecting to dynamo...")
			db = dynamodb.New(session.New(&aws.Config{
				Endpoint:    aws.String(localDynamo),
				Credentials: credentials.NewStaticCredentials("local", "test", "stuff"),
				Region:      aws.String("us-east-1"),
			}))

			if err := initTables(db); err != nil {
				log.WithError(err).Errorln("Failed to init tables!")
			}
		} else {
			db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
		}
	}

	return &Client{
		db:          db,
		blogService: NewBlogService(db, os.Getenv("INKWELL_BLOGS_TABLE")),
	}
}

// BlogService returns a dynamo implementation of the inkwell BlogService interface.
func (c *Client) BlogService() inkwell.BlogService {
	return c.blogService
}

func initTables(db dynamodbiface.DynamoDBAPI) error {
	return initTable(db, "blogs", "blog_id")
}
