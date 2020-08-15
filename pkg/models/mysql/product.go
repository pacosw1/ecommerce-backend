package mysql

import (
	"database/sql"
	"project-z/pkg/models"
)

//ProductModel type which wraps a sql.DB connection pool.
type ProductModel struct {
	DB *sql.DB
}

//Insert serves to add a product to the database
func (m *ProductModel) Insert(p models.Product) (int, error) {

	query := `INSERT INTO products (name, description, stock, price, salePrice, created) VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	//perform query on DB
	result, err := m.DB.Exec(query, p.Name, p.Description, p.Stock, p.Price, p.SalePrice)
	if err != nil {
		return 0, err
	}

	//get last id
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil

}
