package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	"github.com/haziba/theplanner/models/food"
	"github.com/haziba/theplanner/services/food/meal"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf(request.HTTPMethod)
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

	return get(mealService)
}

func post(request events.APIGatewayProxyRequest, mealService meal.MealService) (events.APIGatewayProxyResponse, error) {
	var m models.Meal

	err := json.Unmarshal([]byte(request.Body), &m)

	if err != nil {
		log.Printf("error unmarshalling meal: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	_, err = mealService.CreateMeal(m)
	if err != nil {
		log.Printf("error creating meal: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("error marshalling meal: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body: string(data),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: 200,
	}, nil
}

func get(mealService meal.MealService) (events.APIGatewayProxyResponse, error) {
	meals, err := mealService.GetAllMeals()

	if err != nil {
		log.Printf("error getting meals %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	m := mealResponse{
		Meals: meals,
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
	Meals []models.Meal `json:"meals"`
}

func createMealResponse(meals []models.Meal) (events.APIGatewayProxyResponse, error) {
	resp := mealResponse{
		Meals: meals,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
		StatusCode: http.StatusForbidden,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
