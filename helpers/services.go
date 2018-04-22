package helpers

import (
	"github.com/haziba/theplanner/services/food/ingredient"
	idb "github.com/haziba/theplanner/services/food/ingredient/dynamodb"
	"github.com/haziba/theplanner/services/food/meal"
	mdb "github.com/haziba/theplanner/services/food/meal/dynamodb"
	"github.com/haziba/theplanner/services/money/log"
	moneyLogDb "github.com/haziba/theplanner/services/money/log/dynamodb"
	"github.com/pkg/errors"
)

func CreateMealService() (meal.MealService, error) {
	var service meal.MealService

	db, err := GetDynamoDBHandle()
	if err != nil {
		return service, errors.Wrap(err, "error getting db handle")
	}

	return mdb.NewDynamoDBMealService(db), nil
}

func CreateIngredientService() (ingredient.IngredientService, error) {
	var service ingredient.IngredientService

	db, err := GetDynamoDBHandle()
	if err != nil {
		return service, errors.Wrap(err, "error getting db handle")
	}

	return idb.NewDynamoDBIngredientService(db), nil
}

func CreateMoneyLogService() (money.MoneyLogService, error) {
	var service money.MoneyLogService

	db, err := GetDynamoDBHandle()
	if err != nil {
		return service, errors.Wrap(err, "error getting db handle")
	}

	return moneyLogDb.NewDynamoDBMoneyLogService(db), nil
}
