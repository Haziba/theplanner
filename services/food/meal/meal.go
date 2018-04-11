package meal

import "github.com/haziba/theplanner/models/food"

type MealService interface {
	CreateMeal(models.Meal) (models.Meal, error)
	GetMeal(id string) (*models.Meal, error)
	GetAllMeals() ([]models.Meal, error)
}
