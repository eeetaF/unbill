package models

type ProductUnit struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price"`
}
