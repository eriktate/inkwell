package dynamo

import (
	"fmt"

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

// Write attempts to create or overwrite a blog in dynamo.
func (s BlogService) Write(blog inkwell.Blog) error {
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
func (s BlogService) Revise(blogID, contentLoc string) error {
	rbi, err := s.reviseBlogInput(blogID, contentLoc)
	if err != nil {
		return err
	}

	if _, err := s.db.UpdateItem(rbi); err != nil {
		return err
	}

	return nil
}

// Publish will mark a given blog as published.
func (s BlogService) Publish(blogID string) error {
	pbi, err := s.publishBlogInput(blogID)
	if err != nil {
		return err
	}

	if _, err := s.db.UpdateItem(pbi); err != nil {
		return err
	}

	return nil
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

func (s BlogService) publishBlogInput(blogID string) (*dynamodb.UpdateItemInput, error) {
	attrKey, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"blog_id": attrKey,
	}

	publishVal := &dynamodb.AttributeValue{
		BOOL: aws.Bool(true),
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

func (s BlogService) reviseBlogInput(blogID, contentLoc string) (*dynamodb.UpdateItemInput, error) {
	attrKey, err := dynamodbattribute.Marshal(blogID)
	if err != nil {
		return nil, err
	}

	key := map[string]*dynamodb.AttributeValue{
		"blog_id": attrKey,
	}

	reviseVal := &dynamodb.AttributeValue{
		S: aws.String(contentLoc),
	}

	expressionNames := map[string]*string{
		"#contentLoc": aws.String("content_loc"),
	}

	expressionValues := map[string]*dynamodb.AttributeValue{
		":contentLoc": reviseVal,
	}

	expression := fmt.Sprintf("SET #contentLoc = :contentLoc")

	return &dynamodb.UpdateItemInput{
		Key: key,
		ExpressionAttributeNames:  expressionNames,
		ExpressionAttributeValues: expressionValues,
		UpdateExpression:          aws.String(expression),
		TableName:                 aws.String(s.blogTable),
	}, nil
}
