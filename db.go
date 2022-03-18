package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare new DynamoDB session
var sess = session.Must(session.NewSession())
var db = dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-2"))
var tableName = "Books-API"

func getItem(pk, sk string) (*Book, error) {
	// Input for dynamoDB query must be formatted
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
	}

	// Get item if found or return if empty
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	// result.Item must be formatted back to our Book struct
	book := new(Book)
	err = dynamodbattribute.UnmarshalMap(result.Item, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func putItem(book *Book) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(book.ISBN),
			},
			"sk": {
				S: aws.String(book.Author),
			},
			"title": {
				S: aws.String(book.Title),
			},
			"itemCount": {
				N: aws.String("1"),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

func incrementItem(pk, sk string) (*Book, error) {
	// Input for dynamoDB query must be formatted
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"pk": {
				S: aws.String(pk),
			},
			"sk": {
				S: aws.String(sk),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc": {
				N: aws.String("1"),
			},
		},
		UpdateExpression: aws.String("set itemCount = itemCount + :inc"),
		ReturnValues:     aws.String("ALL_NEW"),
	}

	// Process update and return updated Book
	result, err := db.UpdateItem(input)
	if err != nil {
		return nil, err
	}

	// result.Item must be formatted back to our Book struct
	book := new(Book)
	err = dynamodbattribute.UnmarshalMap(result.Attributes, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}
