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

	log.Printf("Handle request! %v", request.HTTPMethod)

	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			StatusCode: 200,
		}, nil
	}
	log.Printf("Pre-mealService")
	mealService, err := helpers.CreateMealService()
	if err != nil {
		log.Printf("error creating meal service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}
	log.Printf("Pre-ingredientService")
	ingredientService, err := helpers.CreateIngredientService()
	if err != nil {
		log.Printf("error creating ingredient service: %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}
	log.Printf("Pre-getAllMeals")
	meals, err := mealService.GetAllMeals()

	if err != nil {
		log.Printf("error getting meals %v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	log.Printf("Got services")

	ingredients, err := ingredientService.GetAllIngredients()

	if err != nil {
		log.Printf("error getting ingredients%v\n", err)
		return helpers.CreateInternalServerErrorResponse()
	}

	log.Printf("Got incredients")

	var plannerRequestBody plannerRequestBody
	err = json.Unmarshal([]byte(request.Body), &plannerRequestBody)
	if err != nil {
		log.Printf("error unmarshalling `%v`\n", request.Body)
		return events.APIGatewayProxyResponse{}, err
	}

	log.Printf("Got body")

	var chosenMeals []models.Meal
	var unchosenMeals []models.Meal

	for _, meal := range meals {
		chosen := false

		for _, preChosenMeal := range plannerRequestBody.ChosenMeals {
			if meal.Id == preChosenMeal.MealID {
				chosen = true
			}
		}

		if chosen {
			chosenMeals = append(chosenMeals, meal)
		} else {
			unchosenMeals = append(unchosenMeals, meal)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(unchosenMeals), func(i, j int) { unchosenMeals[i], unchosenMeals[j] = unchosenMeals[j], unchosenMeals[i] })

	for i := len(chosenMeals); i < 4 && len(unchosenMeals) > 0; i++ {
		unchosenID := rand.Intn(len(unchosenMeals))
		chosenMeals = append(chosenMeals, unchosenMeals[unchosenID])
		unchosenMeals = append(unchosenMeals[:unchosenID], unchosenMeals[unchosenID+1:]...)
	}

	log.Printf("Got meals")

	m := mealPlannerResponse{
		Meals:       chosenMeals,
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

type plannerRequestBody struct {
	ChosenMeals []struct {
		MealID              string `json:"mealId"`
		OptionalIngredients []struct {
			Ingredient models.Ingredient `json:"ingredient"`
			Use        bool              `json:"use"`
		} `json:"optionalIngredients"`
	} `json:"chosenMeals"`
}

type mealPlannerResponse struct {
	Meals       []models.Meal       `json:"meals"`
	Ingredients []models.Ingredient `json:"ingredients"`
}

type mealIngredient struct {
	Ingredient []models.Ingredient
	Quantity   int
}

func main() {
	log.Printf("Booting...")
	lambda.Start(handleRequest)
}
