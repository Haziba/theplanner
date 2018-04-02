package helpers

import (
	"github.com/haziba/theplanner/services/meal"
	mdb "github.com/haziba/theplanner/services/meal/dynamodb"
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
