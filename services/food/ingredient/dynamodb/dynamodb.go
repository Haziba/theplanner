package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haziba/theplanner/models/food"
	"github.com/pkg/errors"
	"github.com/satori/uuid"
)

const tableName = "tp-ingredients"

type DynamoDBIngredientService struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBIngredientService(db *dynamodb.DynamoDB) DynamoDBIngredientService {
	return DynamoDBIngredientService{
		db: db,
	}
}

func (s DynamoDBIngredientService) CreateIngredient(i models.Ingredient) (models.Ingredient, error) {
	if i.Id == "" {
		i.Id = uuid.NewV4().String()
	}

	item, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return i, errors.Wrap(err, "error marshalling ingredient")
	}

	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return i, errors.Wrap(err, "error putting ingredient")
	}

	return i, nil
}

func (s DynamoDBIngredientService) GetAllIngredients() ([]models.Ingredient, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := s.db.Scan(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting ingredients")
	}

	var i []models.Ingredient
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &i)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling items")
	}

	return i, nil
}
