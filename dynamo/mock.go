package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamo struct {
	dynamodbiface.DynamoDBAPI

	items map[string]map[string]*dynamodb.AttributeValue

	PutItemCalled int
}

func NewMockDynamo() *MockDynamo {
	return &MockDynamo{items: make(map[string]map[string]*dynamodb.AttributeValue)}
}

func (m *MockDynamo) AddItem(tableName string, item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	m.items[tableName] = av
	return nil
}

func (m *MockDynamo) GetItem(gii *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	var gio dynamodb.GetItemOutput
	gio.Item = m.items[*gii.TableName]

	return &gio, nil
}

func (m *MockDynamo) PutItem(pii *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.PutItemCalled += 1
	return &dynamodb.PutItemOutput{}, nil
}
