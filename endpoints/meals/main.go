package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	"github.com/haziba/theplanner/models"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	mealService, err := helpers.CreateMealService()
	if err != nil {
		log.Printf("error creating meal service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

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
		Body:       string(data),
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
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
