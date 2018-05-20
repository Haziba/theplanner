package models

type Meal struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Recipe     string `json:"recipe"`
	TimeToMake struct {
		Time       int  `json:"time"`
		Confidence int  `json:"confidence"`
		HasABreak  bool `json:"has-a-break"`
	} `json:"time-to-make"`
	Ingredients []MealIngredient `json:"ingredients"`
}

type MealIngredient struct {
	Id          string  `json:"id"`
	Quantity    float32 `json:"quantity"`
	Measurement string  `json:"measurement"`
	Optional    bool    `json:"optional"`
}
