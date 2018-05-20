package ingredient

import "github.com/haziba/theplanner/models/food"

type IngredientService interface {
	CreateIngredient(models.Ingredient) (models.Ingredient, error)
	UpdateIngredient(models.Ingredient) (models.Ingredient, error)
	GetAllIngredients() ([]models.Ingredient, error)
	GetIngredient(id string) (*models.Ingredient, error)
}
