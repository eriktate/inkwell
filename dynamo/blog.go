package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/eriktate/inkwell"
)

type BlogService struct {
	blogTable string
	db        dynamodbiface.DynamoDBAPI
}

func NewUserService(db dynamodbiface.DynamoDBAPI, blogTable string) BlogService {
	return BlogService{
		blogTable: blogTable,
		db:        db,
	}
}

func (s *BlogService) Get(blogID string) (inkwell.Blog, error) {
	s.db.GetItem
}

func getBlogInput(blogID string) inkwell.Blog {

}
