package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	"github.com/haziba/theplanner/models/food"
	"github.com/haziba/theplanner/services/food/meal"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			StatusCode: 200,
		}, nil
	}

	mealService, err := helpers.CreateMealService()
	if err != nil {
		log.Printf("error creating meal service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if request.HTTPMethod == "POST" {
		return post(request, mealService)
	}

	return get(mealService, request.PathParameters["meal-id"])
}

func post(request events.APIGatewayProxyRequest, mealService meal.MealService) (events.APIGatewayProxyResponse, error) {
	var m []models.Meal

	err := json.Unmarshal([]byte(request.Body), &m)

	if err != nil {
		log.Printf("error unmarshalling meal: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	for _, meal := range m {
		meal, err = mealService.CreateMeal(meal)
		if err != nil {
			log.Printf("error creating meal: %v\n", err)
			return helpers.CreateBadRequestResponse()
		}
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("error marshalling meal: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: 200,
	}, nil
}

func get(mealService meal.MealService, mealID string) (events.APIGatewayProxyResponse, error) {
	meal, err := mealService.GetMeal(mealID)

	if err != nil {
		log.Printf("error getting meals %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if meal == nil {
		log.Printf("couldn't find meal %v\n", mealID)
		return helpers.CreateNotFoundResponse()
	}

	ingredientService, err := helpers.CreateIngredientService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if err != nil {
		log.Printf("error getting ingredients%v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	meal, ingredients, err := helpers.GetMealAndIngredients(mealService, ingredientService, mealID)

	if err != nil {
		log.Printf("error getting meal and ingredients: %v\n", err)
		return helpers.CreateNotFoundResponse()
	}

	m := mealResponse{
		Meal:        *meal,
		Ingredients: *ingredients,
	}

	data, err := json.Marshal(m)
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

type mealResponse struct {
	Meal        models.Meal         `json:"meal"`
	Ingredients []models.Ingredient `json:"ingredients"`
}

func main() {
	lambda.Start(handleRequest)
}
