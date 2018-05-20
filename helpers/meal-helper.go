package helpers

import (
	"log"

	models "github.com/haziba/theplanner/models/food"
	"github.com/haziba/theplanner/services/food/ingredient"
	"github.com/haziba/theplanner/services/food/meal"
)

func GetMealAndIngredients(mealService meal.MealService, ingredientService ingredient.IngredientService, mealID string) (*models.Meal, *[]models.Ingredient, error) {
	meal, err := mealService.GetMeal(mealID)

	if err != nil {
		log.Printf("error getting meals %v\n", err)
		return nil, nil, err
	}

	if meal == nil {
		log.Printf("couldn't find meal %v\n", mealID)
		return nil, nil, err
	}

	ingredients, err := ingredientService.GetAllIngredients()

	return meal, &ingredients, nil
}
