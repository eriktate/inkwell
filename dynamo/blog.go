package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/eriktate/inkwell"
	log "github.com/sirupsen/logrus"
)

type BlogService struct {
	blogTable string
	db        dynamodbiface.DynamoDBAPI
}

func NewBlogService(db dynamodbiface.DynamoDBAPI, blogTable string) BlogService {
	return BlogService{
		blogTable: blogTable,
		db:        db,
	}
}

func (s BlogService) Get(blogID string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	gbi, err := s.getBlogInput(blogID)
	if err != nil {
		log.WithError(err).Error("Failed to create GetItemInput.")
		return blog, err
	}

	gio, err := s.db.GetItem(gbi)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve blog from dynamo.")
		return blog, err
	}

	if err := dynamodbattribute.UnmarshalMap(gio.Item, &blog); err != nil {
		log.WithError(err).Error("Failed to unmarshal blog from dynamo.")
		return blog, err
	}

	return blog, nil
}

func (s BlogService) Write(blog inkwell.Blog) error {
	pbi, err := s.putBlogInput(blog)
	if err != nil {
		return err
	}

	// TODO: Maybe do something with the output at some point?
	_, err = s.db.PutItem(pbi)
	return err
}

func (s BlogService) getBlogInput(blogID string) (*dynamodb.GetItemInput, error) {
	av, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"blog_id": av,
	}

	return &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(s.blogTable),
	}, nil
}

func (s BlogService) putBlogInput(blog inkwell.Blog) (*dynamodb.PutItemInput, error) {
	avMap, err := dynamodbattribute.MarshalMap(blog)
	if err != nil {
		return nil, err
	}

	return &dynamodb.PutItemInput{
		Item:      avMap,
		TableName: aws.String(s.blogTable),
	}, nil
}
