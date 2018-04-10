package models

type Meal struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	TimeToMake struct {
		Time       int  `json: "time"`
		Confidence int  `json: "confidence"`
		HasABreak  bool `json:"has-a-break"`
	} `json:"time-to-make"`
	Ingredients []struct {
		Id       string `json:"id"`
		Quantity int    `json:"quantity"`
		Optional bool   `json:"optional"`
	} `json:"ingredients"`
}
