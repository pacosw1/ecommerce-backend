package main

import "project-z/pkg/models/mysql"

//Models wrapper to group all model dependencies in one place
type Models struct {
	Product  *mysql.ProductModel
	Category *mysql.CategoryModel
}
