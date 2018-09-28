package models

type Ingredient struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Tesco *struct {
		QuantityAsSold int     `json:"quantityAsSold"`
		TPNC           string  `json:"tpnc"`
		Price          float32 `json:"price"`
	} `json:"tesco"`
}
