package dynamodb

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haziba/theplanner/models/food"
	"github.com/pkg/errors"
	"github.com/satori/uuid"
)

const tableName = "tp-food-planner"

type DynamoDBPlannerService struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBPlannerService(db *dynamodb.DynamoDB) DynamoDBPlannerService {
	return DynamoDBPlannerService{
		db: db,
	}
}

func (s DynamoDBPlannerService) GetPlanner(id string) (*models.Planner, error) {
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

	var p models.Planner
	dynamodbattribute.UnmarshalMap(output.Item, &p)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling item")
	}

	return &p, nil
}

func (s DynamoDBPlannerService) CreatePlanner(p models.Planner) (models.Planner, error) {
	p.ID = uuid.NewV4().String()
	p.When = time.Now()

	item, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return p, errors.Wrap(err, "error marshalling planner")
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return p, errors.Wrap(err, "error putting planner")
	}

	return p, nil
}

func (s DynamoDBPlannerService) GetAllPlanners() ([]models.Planner, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := s.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting planners")
	}

	var p []models.Planner
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &p)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling items")
	}

	return p, nil
}
