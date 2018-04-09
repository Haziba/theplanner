package ingredient

import "github.com/haziba/theplanner/models"

type IngredientService interface {
	CreateIngredient(models.Ingredient) (models.Ingredient, error)
	GetAllIngredients() ([]models.Ingredient, error)
}
