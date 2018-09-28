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
	"github.com/haziba/theplanner/services/food/ingredient"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Holla holla!")

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

	ingredientService, err := helpers.CreateIngredientService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	if request.HTTPMethod == "POST" {
		return post(request, ingredientService)
	}

	return get(ingredientService)
}

func post(request events.APIGatewayProxyRequest, ingredientService ingredient.IngredientService) (events.APIGatewayProxyResponse, error) {
	var i models.Ingredient

	err := json.Unmarshal([]byte(request.Body), &i)

	if err != nil {
		log.Printf("error unmarshalling ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	_, err = ingredientService.CreateIngredient(i)
	if err != nil {
		log.Printf("error creating ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	data, err := json.Marshal(i)
	if err != nil {
		log.Printf("error marshalling ingredient: %v\n", err)
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

func get(ingredientService ingredient.IngredientService) (events.APIGatewayProxyResponse, error) {
	ingredients, err := ingredientService.GetAllIngredients()

	if err != nil {
		log.Printf("error getting ingredients %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	m := ingredientResponse{
		Ingredients: ingredients,
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

type ingredientResponse struct {
	Ingredients []models.Ingredient `json:"ingredients"`
}

func createIngredientResponse(ingredients []models.Ingredient) (events.APIGatewayProxyResponse, error) {
	resp := ingredientResponse{
		Ingredients: ingredients,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body: string(data),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
