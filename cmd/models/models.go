package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

//Product stores product data
type Product struct {
	ID          string    `validate:"omitempty"`
	Name        string    `validate:"required" ` //required
	Stock       uint16    `validate:"required"`  //required
	Description string    `validate:"required"`  //required
	Price       float32   `validate:"required"`  //required
	SalePrice   float32   `validate:"omitempty"`
	Created     time.Time `validate:"omitempty"`
	Images      []string  `validate:"-"`
	Tags        []string  `json:"-"`
}

//Specs serves to calculate shipping price and package dimensions
type Specs struct {
	Weigth float32
	Length float32
	Width  float32
	Height float32
	Unit   string
}

//Category serves to organize products
type Category struct {
	Name        string
	Description string
	ParentID    string
}

// //Order store purchase and shipping info
// type Order struct {
// 	Date     uint64
// 	Subtotal float32
// 	Total    float32
// 	items    []*Product
// 	Status   string
// 	Shipping *Shipping
// 	UserID   string
// }

//Address serves to store user shipping address
type Address struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Col        string
	Country    string
}
