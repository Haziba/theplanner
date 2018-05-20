package models

import "time"

type Planner struct {
	ID          string              `json:"id"`
	When        time.Time           `json:"when"`
	ChosenMeals []PlannerChosenMeal `json:"chosenMeal"`
}

type PlannerChosenMeal struct {
	MealID      string   `json:"mealId"`
	Ingredients []string `json:"ingredients"`
}
