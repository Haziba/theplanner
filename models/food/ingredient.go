package models

import money "github.com/rhymond/go-money"

type Ingredient struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Tesco *struct {
		QuantityAsSold int         `json:"quantityAsSold"`
		TPNC           string      `json:"tpnc"`
		Price          money.Money `json:"price"`
	} `json:"tesco"`
}
