package mysql

import (
	"database/sql"
	"project-z/pkg/models"
)

//CategoryModel type which wraps a sql.DB connection pool.
type CategoryModel struct {
	DB *sql.DB
}

// //Insert serves to add a product to the database
// func (m *CategoryModel) Insert(title, descp) (int, error) {
// 	return 0, nil
// }

//Get serves to get a product give an id parameter
func (m *CategoryModel) Get(id int) (*models.Category, error) {
	return nil, nil
}
