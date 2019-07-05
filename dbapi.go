package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	partitionKey = "ID"
)

func tableStruct(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(partitionKey),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(partitionKey),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest), // ondemand
		TableName:   aws.String(tableName),
	}
}

func validTable(db dynamodb.DynamoDB, tableName string) bool {
	input := &dynamodb.ListTablesInput{}
	isTableName := false
	// 한번에 최대 100개의 테이블만 가지고 올 수 있다.
	// 한 리전에 최대 256개의 테이블이 존재할 수 있다.
	// https://docs.aws.amazon.com/ko_kr/amazondynamodb/latest/developerguide/Limits.html
	for {
		result, err := db.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Fprintf(os.Stderr, "%s %s\n", dynamodb.ErrCodeInternalServerError, err.Error())
				default:
					fmt.Fprintf(os.Stderr, "%s\n", aerr.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			return false
		}

		for _, n := range result.TableNames {
			if *n == tableName {
				isTableName = true
				break
			}
		}
		if isTableName {
			break
		}
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
	return isTableName
}

func hasItem(db dynamodb.DynamoDB, tableName string, primarykey string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(primarykey),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return false, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

// AddCharacter 는 사용자를 추가하는 함수이다.
func AddCharacter(db dynamodb.DynamoDB) error {
	hasBool, err := hasItem(db, *flagTable, *flagID)
	if err != nil {
		return err
	}
	if hasBool {
		return errors.New("The data already exists. Can not add data")
	}

	c := Character{
		ID:              *flagID,
		Regnum:          *flagRegnum,
		Manager:         *flagManager,
		FieldOfActivity: *flagFieldOfActivity,
		Concept:         *flagConcept,
		StartDate:       *flagStartDate,
		Email:           *flagEmail,
	}

	dynamodbJSON, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}

	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// SetCharacter 는 프로젝트 자료구조를 수정하는 함수이다.
func SetCharacter(db dynamodb.DynamoDB) error {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(*flagTable),
		Key: map[string]*dynamodb.AttributeValue{
			partitionKey: {
				S: aws.String(*flagID),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return err
	}
	c := Character{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &c)
	if err != nil {
		return err
	}

	if *flagRegnum != "" && c.Regnum != *flagRegnum {
		c.Regnum = *flagRegnum
	}
	if *flagManager != "" && c.Manager != *flagManager {
		c.Manager = *flagManager
	}
	if *flagFieldOfActivity != "" && c.FieldOfActivity != *flagFieldOfActivity {
		c.FieldOfActivity = *flagFieldOfActivity
	}
	if *flagConcept != "" && c.Concept != *flagConcept {
		c.Concept = *flagConcept
	}
	if *flagStartDate != "" && c.StartDate != *flagStartDate {
		c.StartDate = *flagStartDate
	}
	if *flagEmail != "" && c.Email != *flagEmail {
		c.Email = *flagEmail
	}

	dynamodbJSON, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return err
	}
	data := &dynamodb.PutItemInput{
		Item:      dynamodbJSON,
		TableName: aws.String(*flagTable),
	}
	_, err = db.PutItem(data)
	if err != nil {
		return err
	}
	return nil
}

// GetCharacters 는 사용자를 가지고오는 함수이다.
func GetCharacters(db dynamodb.DynamoDB, word string) error {
	proj := expression.NamesList(
		expression.Name("ID"),
		expression.Name("Regnum"),
		expression.Name("Manager"),
		expression.Name("FieldOfActivity"),
		expression.Name("Concept"),
		expression.Name("StartDate"),
		expression.Name("Email"),
	)

	f1 := expression.Name("ID").Contains(word)
	f2 := expression.Name("Regnum").Contains(word)
	f3 := expression.Name("Manager").Contains(word)
	f4 := expression.Name("FieldOfActivity").Contains(word)
	f5 := expression.Name("Concept").Contains(word)
	f6 := expression.Name("StartDate").Contains(word)
	f7 := expression.Name("Email").Contains(word)

	expr, err := expression.NewBuilder().
		WithFilter(f1.Or(f2).Or(f3).Or(f4).Or(f5).Or(f6).Or(f7)).
		WithProjection(proj).
		Build()
	if err != nil {
		return err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(*flagTable),
	}
	result, err := db.Scan(params)
	if err != nil {
		return err
	}
	for _, i := range result.Items {
		c := Character{}
		err = dynamodbattribute.UnmarshalMap(i, &c)
		if err != nil {
			return err
		}
		fmt.Println(c)
	}
	return nil
}
