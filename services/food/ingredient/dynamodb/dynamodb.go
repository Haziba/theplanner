package dynamodb

import (
	"fmt"

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

func (s DynamoDBIngredientService) GetIngredient(id string) (*models.Ingredient, error) {
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
		return nil, errors.Wrap(err, "error getting ingredient")
	}
	if output.Item == nil {
		return nil, nil
	}

	var i models.Ingredient
	dynamodbattribute.UnmarshalMap(output.Item, &i)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling item")
	}

	return &i, nil
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

func (s DynamoDBIngredientService) UpdateIngredient(i models.Ingredient) (models.Ingredient, error) {
	item, err := dynamodbattribute.MarshalMap(i)
	fmt.Printf("Map marshalled")
	if err != nil {
		return i, errors.Wrap(err, "error marshalling ingredient")
	}
	fmt.Printf("PutItem %v", tableName)
	_, err = s.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	fmt.Printf("Post put")
	if err != nil {
		return i, errors.Wrap(err, "error putting ingredient")
	}

	return i, nil
}

func (s DynamoDBIngredientService) GetAllIngredients() ([]models.Ingredient, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	fmt.Printf("Table name %v\n", tableName)
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
