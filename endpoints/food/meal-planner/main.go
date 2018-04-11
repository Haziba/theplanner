package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	helpers "github.com/haziba/theplanner/helpers"
	"github.com/haziba/theplanner/models/food"
)

func handleRequest(context context.Context,
	request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	meals, err := mealService.GetAllMeals()

	if err != nil {
		log.Printf("error getting meals %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	ingredients, err := ingredientService.GetAllIngredients()

	if err != nil {
		log.Printf("error getting ingredients%v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(meals), func(i, j int) { meals[i], meals[j] = meals[j], meals[i] })

	chosenMeals := meals[:4]
	var mealNames []string
	mealIngredients := make(map[string]int)

	for _, meal := range chosenMeals {
		mealNames = append(mealNames, meal.Name)

		for _, mealIngredient := range meal.Ingredients {
			var ingredientName string
			for _, ingredient := range ingredients {
				if ingredient.Id == mealIngredient.Id {
					ingredientName = ingredient.Name
					break
				}
			}

			if ingredientName == "" {
				log.Printf("cannot find ingredient %v\n", mealIngredient.Id)
			}

			mealIngredients[ingredientName] = mealIngredients[ingredientName] + mealIngredient.Quantity
		}
	}

	m := mealPlannerResponse{
		Meals:       mealNames,
		Ingredients: mealIngredients,
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

type mealPlannerResponse struct {
	Meals       []string       `json:"meals"`
	Ingredients map[string]int `json:"ingredients"`
}

type mealIngredient struct {
	Ingredient []models.Ingredient
	Quantity   int
}

func main() {
	lambda.Start(handleRequest)
}
