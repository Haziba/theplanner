package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haziba/theplanner/models/food"
	"github.com/pkg/errors"
	"github.com/satori/uuid"
)

const tableName = "tp-meals"

type DynamoDBMealService struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBMealService(db *dynamodb.DynamoDB) DynamoDBMealService {
	return DynamoDBMealService{
		db: db,
	}
}

func (s DynamoDBMealService) GetMeal(id string) (*models.Meal, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}

	output, err := s.db.GetItem(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting meal")
	}
	if output.Item == nil {
		return nil, nil
	}

	var m models.Meal
	dynamodbattribute.UnmarshalMap(output.Item, &m)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling item")
	}

	return &m, nil
}

func (s DynamoDBMealService) CreateMeal(m models.Meal) (models.Meal, error) {
	m.Id = uuid.NewV4().String()

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

func (s DynamoDBMealService) GetAllMeals() ([]models.Meal, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := s.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting meals")
	}

	var m []models.Meal
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &m)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling items")
	}

	return m, nil
}
