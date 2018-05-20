package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/haziba/theplanner/helpers"
	models "github.com/haziba/theplanner/models/food"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			StatusCode: 200,
		}, nil
	}

	mealService, err := helpers.CreateMealService()
	if err != nil {
		log.Printf("error creating meal service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	ingredientService, err := helpers.CreateIngredientService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	var plannerRequestBody plannerRequestBody
	err = json.Unmarshal([]byte(request.Body), &plannerRequestBody)
	if err != nil {
		log.Printf("error unmarshalling `%v`\n", request.Body)
		return events.APIGatewayProxyResponse{}, err
	}

	ingredients, err := ingredientService.GetAllIngredients()

	var planner models.Planner
	var plannerIngredients []ChosenIngredient

	for _, chosenMeal := range plannerRequestBody.ChosenMeals {
		var plannerChosenMeal models.PlannerChosenMeal
		meal, _ := mealService.GetMeal(chosenMeal.MealID)

		for _, ingredient := range meal.Ingredients {
			if !ingredient.Optional || contains(chosenMeal.OptionalIngredients, ingredient.Id) {
				plannerChosenMeal.Ingredients = append(plannerChosenMeal.Ingredients, ingredient.Id)
				plannerIngredients = addToPlannerIngredients(plannerIngredients, ingredients, ingredient)
			}
		}

		planner.ChosenMeals = append(planner.ChosenMeals, plannerChosenMeal)
	}

	plannerService, err := helpers.CreatePlannerService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	plannerService.CreatePlanner(planner)

	data, err := json.Marshal(plannerIngredients)
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

func contains(s []OptionalIngredient, ingredientID string) bool {
	for _, a := range s {
		if a.Ingredient.Id == ingredientID {
			return true
		}
	}
	return false
}

func find(s []models.Ingredient, ingredientID string) *models.Ingredient {
	for _, a := range s {
		if a.Id == ingredientID {
			return &a
		}
	}

	return nil
}

func addToPlannerIngredients(s []ChosenIngredient, ingredients []models.Ingredient, ingred models.MealIngredient) []ChosenIngredient {
	for _, a := range s {
		if a.Ingredient.Id == ingred.Id {
			a.Quantity = a.Quantity + ingred.Quantity
			return s
		}
	}

	for _, a := range ingredients {
		if a.Id == ingred.Id {
			return append(s, ChosenIngredient{a, ingred.Quantity})
		}
	}

	return s
}

type plannerRequestBody struct {
	ChosenMeals []struct {
		MealID              string               `json:"mealId"`
		OptionalIngredients []OptionalIngredient `json:"optionalIngredients"`
	} `json:"chosenMeals"`
}

type OptionalIngredient struct {
	Ingredient models.Ingredient `json:"ingredient"`
	Use        bool              `json:"use"`
}

type ChosenIngredient struct {
	Ingredient models.Ingredient `json:"ingredient"`
	Quantity   float32           `json:"quantity"`
}

func main() {
	lambda.Start(handleRequest)
}
