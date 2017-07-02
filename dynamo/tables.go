package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func initTable(db dynamodbiface.DynamoDBAPI, tableName string, partitionKey string) error {
	attributeDef := &dynamodb.AttributeDefinition{
		AttributeName: aws.String(partitionKey),
		AttributeType: aws.String("S"),
	}

	attributeDefs := []*dynamodb.AttributeDefinition{attributeDef}

	keyDef := &dynamodb.KeySchemaElement{
		AttributeName: aws.String(partitionKey),
		KeyType:       aws.String("HASH"),
	}

	keySchema := []*dynamodb.KeySchemaElement{keyDef}

	throughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(20),
		WriteCapacityUnits: aws.Int64(10),
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefs,
		KeySchema:             keySchema,
		TableName:             aws.String(tableName),
		ProvisionedThroughput: throughput,
	}

	_, err := db.CreateTable(createTableInput)
	return err
}
