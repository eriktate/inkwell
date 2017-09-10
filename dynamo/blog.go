package dynamo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/eriktate/inkwell"
	log "github.com/sirupsen/logrus"
)

// BlogService is a dynamo implementatino of the inkwell BlogService interface.
type BlogService struct {
	blogTable string
	db        dynamodbiface.DynamoDBAPI
}

// NewBlogService returns a dynamo implementation of the inkwell BlogService
// interface.
func NewBlogService(db dynamodbiface.DynamoDBAPI, blogTable string) BlogService {
	return BlogService{
		blogTable: blogTable,
		db:        db,
	}
}

// Get returns data related to a given blog that's stored in dynamo.
func (s BlogService) Get(authorID, blogID string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	gbi, err := s.getBlogInput(authorID, blogID)
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

// Write attempts to create or overwrite a blog in dynamo.
func (s BlogService) Write(blog inkwell.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	// check for existing
	existing, err := s.Get(blog.AuthorID, blog.ID)
	if err != nil {
		return err
	}

	// if existing, don't change created at date
	if existing.ID != "" {
		blog.CreatedAt = existing.CreatedAt
	}

	pbi, err := s.putBlogInput(blog)
	if err != nil {
		return err
	}

	// TODO: Maybe do something with the output at some point?
	_, err = s.db.PutItem(pbi)
	return err
}

// Revise will update the content for a given blog. In dynamo, this means
// updating the content location (likely housed in s3). This is an action that
// probably won't be performed very often.
func (s BlogService) Revise(authorID, blogID, contentLoc string) error {
	rbi, err := s.reviseBlogInput(authorID, blogID, contentLoc)
	if err != nil {
		return err
	}

	if _, err := s.db.UpdateItem(rbi); err != nil {
		return err
	}

	return nil
}

// Publish will mark a given blog as published.
func (s BlogService) Publish(authorID, blogID string) error {
	pbi, err := s.publishBlogInput(authorID, blogID, true)
	if err != nil {
		return err
	}

	if _, err := s.db.UpdateItem(pbi); err != nil {
		return err
	}

	return nil
}

// Redact will mark a given blog as unpublished.
func (s BlogService) Redact(authorID, blogID string) error {
	rbi, err := s.publishBlogInput(authorID, blogID, false)
	if err != nil {
		return err
	}

	if _, err := s.db.UpdateItem(rbi); err != nil {
		return err
	}

	return nil
}

func (s BlogService) getBlogInput(authorID, blogID string) (*dynamodb.GetItemInput, error) {
	authorAV, err := dynamodbattribute.Marshal(authorID)
	if err != nil {
		return nil, err
	}

	blogAV, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"author_id": authorAV,
		"blog_id":   blogAV,
	}

	return &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(s.blogTable),
	}, nil
}

func (s BlogService) putBlogInput(blog inkwell.Blog) (*dynamodb.PutItemInput, error) {
	avMap, err := dynamodbattribute.MarshalMap(&blog)
	if err != nil {
		return nil, err
	}

	log.WithField("avMap", avMap).Println("About to write")
	return &dynamodb.PutItemInput{
		Item:      avMap,
		TableName: aws.String(s.blogTable),
	}, nil
}

func (s BlogService) publishBlogInput(authorID, blogID string, status bool) (*dynamodb.UpdateItemInput, error) {
	authorAV, err := dynamodbattribute.Marshal(authorID)
	if err != nil {
		return nil, err
	}

	blogAV, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"author_id": authorAV,
		"blog_id":   blogAV,
	}

	publishVal := &dynamodb.AttributeValue{
		BOOL: aws.Bool(status),
	}

	expressionNames := map[string]*string{
		"#published": aws.String("published"),
	}

	expressionValues := map[string]*dynamodb.AttributeValue{
		":published": publishVal,
	}

	expression := fmt.Sprintf("SET #published = :published")

	return &dynamodb.UpdateItemInput{
		Key: key,
		ExpressionAttributeNames:  expressionNames,
		ExpressionAttributeValues: expressionValues,
		UpdateExpression:          aws.String(expression),
		TableName:                 aws.String(s.blogTable),
	}, nil
}

func (s BlogService) reviseBlogInput(authorID, blogID, contentLoc string) (*dynamodb.UpdateItemInput, error) {
	authorAV, err := dynamodbattribute.Marshal(authorID)
	if err != nil {
		return nil, err
	}

	blogAV, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"author_id": authorAV,
		"blog_id":   blogAV,
	}

	now := time.Now()
	updatedAt, err := dynamodbattribute.Marshal(&now)
	if err != nil {
		return nil, err
	}

	expressionNames := map[string]*string{
		"#updatedAt": aws.String("updated_at"),
	}

	expressionValues := map[string]*dynamodb.AttributeValue{
		":updatedAt": updatedAt,
	}
	// expressionNames := map[string]*string{
	// 	"#contentLoc": aws.String("content_loc"),
	// }

	// expressionValues := map[string]*dynamodb.AttributeValue{
	// 	":contentLoc": reviseVal,
	// }

	// expression := fmt.Sprintf("SET #contentLoc = :contentLoc")
	expression := fmt.Sprintf("SET #updatedAt = :updatedAt")

	return &dynamodb.UpdateItemInput{
		Key: key,
		ExpressionAttributeNames:  expressionNames,
		ExpressionAttributeValues: expressionValues,
		UpdateExpression:          aws.String(expression),
		TableName:                 aws.String(s.blogTable),
	}, nil
}
