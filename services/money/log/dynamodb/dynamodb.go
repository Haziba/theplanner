package dynamodb

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haziba/theplanner/models/money"
	"github.com/pkg/errors"
	"github.com/satori/uuid"
)

const tableName = "tp-money-logs"

type DynamoDBMoneyLogService struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBMoneyLogService(db *dynamodb.DynamoDB) DynamoDBMoneyLogService {
	return DynamoDBMoneyLogService{
		db: db,
	}
}

func (s DynamoDBMoneyLogService) CreateMoneyLog(m models.MoneyLog) (models.MoneyLog, error) {
	m.ID = uuid.NewV4().String()
	m.When = time.Now()

	item, err := dynamodbattribute.MarshalMap(m)
	if err != nil {
		return m, errors.Wrap(err, "error marshalling meal")
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return m, errors.Wrap(err, "error putting meal")
	}

	return m, nil
}

func (s DynamoDBMoneyLogService) GetAllMoneyLogs() ([]models.MoneyLog, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := s.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting money logs")
	}

	var m []models.MoneyLog
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &m)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling items")
	}

	return m, nil
}
