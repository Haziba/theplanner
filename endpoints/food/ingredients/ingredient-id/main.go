package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	models "github.com/haziba/theplanner/models/food"
	"github.com/haziba/theplanner/services/food/ingredient"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "PUT",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			StatusCode: 200,
		}, nil
	}

	ingredientService, err := helpers.CreateIngredientService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if request.HTTPMethod == "PUT" {
		return put(request, ingredientService)
	}

	return get(ingredientService, request.PathParameters["ingredient-id"])
}

func put(request events.APIGatewayProxyRequest, ingredientService ingredient.IngredientService) (events.APIGatewayProxyResponse, error) {
	var i models.Ingredient

	err := json.Unmarshal([]byte(request.Body), &i)

	if err != nil {
		log.Printf("error unmarshalling ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	ingredientService.UpdateIngredient(i)
	log.Printf("Bblawkelkweutts butts")

	return events.APIGatewayProxyResponse{
		//Body:       string(data),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: 200,
	}, nil
}

func get(ingredientService ingredient.IngredientService, ingredientID string) (events.APIGatewayProxyResponse, error) {
	ingredient, err := ingredientService.GetIngredient(ingredientID)

	if err != nil {
		log.Printf("error getting ingredients %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if ingredient == nil {
		log.Printf("couldn't find ingredient %v\n", ingredientID)
		return helpers.CreateNotFoundResponse()
	}

	data, err := json.Marshal(ingredient)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body: string(data),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
