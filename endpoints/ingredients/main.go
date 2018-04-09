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
	"github.com/haziba/theplanner/services/ingredient"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
	var m models.Ingredient

	err := json.Unmarshal([]byte(request.Body), &m)

	if err != nil {
		log.Printf("error unmarshalling ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	m, err = ingredientService.CreateIngredient(m)
	if err != nil {
		log.Printf("error creating ingredient: %v\n", err)
		return helpers.CreateBadRequestResponse()
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("error marshalling ingredient: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       string(data),
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
		Body:       string(data),
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
		Body:       string(data),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}