package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)

	attributeDefinitions := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("Year"),
			AttributeType: aws.String("N"),
		},
		{
			AttributeName: aws.String("Title"),
			AttributeType: aws.String("S"),
		},
	}

	keySchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("Year"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("Title"),
			KeyType:       aws.String("RANGE"),
		},
	}

	provisionedThroughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}

	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchema,
		ProvisionedThroughput: provisionedThroughput,
		TableName:             aws.String("testingTable"),
	})

	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err)
		os.Exit(1)
	}
}
