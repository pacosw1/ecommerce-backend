package mysql

import (
	"database/sql"
	"mime/multipart"
	saver "project-z/pkg/image-saver"
	"project-z/pkg/models"
)

//ProductModel type which wraps a sql.DB connection pool.
type ProductModel struct {
	DB *sql.DB
}

//Insert serves to add a product to the database
func (m *ProductModel) Insert(p models.Product, images []*multipart.FileHeader) error {

	//start transaction
	transaction, err := m.DB.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}

	query := `INSERT INTO products 
			(name, description, stock, price, salePrice, created) 
			VALUES
			(?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	//perform query on DB
	result, err := transaction.Exec(query, p.Name, p.Description, p.Stock, p.Price, p.SalePrice)
	if err != nil {
		transaction.Rollback()
		return err
	}

	//if db saved, get its generate ID
	productID, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}

	//try to save images to disk and return paths to store in DB
	paths, err := saver.SaveImagesToDisk("cmd/static/images", images)

	if err != nil {
		transaction.Rollback()
		//removes images stored since transaction failed
		saver.CleanUp(paths)
		return err
	}

	var values []interface{}
	imageQuery := `INSERT INTO images (productId, path) VALUES `

	for index, path := range paths {
		row := "(?, ?)"
		if index < len(paths)-1 {
			row += ", "
		}
		imageQuery += row
		values = append(values, productID, path)
	}

	//execute query
	result, err = transaction.Exec(imageQuery, values...)

	if err != nil {
		//removes images stored since transaction failed
		saver.CleanUp(paths)
		transaction.Rollback()
		return err
	}

	//we have added everything and created our product so commit transaction and apply changes
	err = transaction.Commit()

	if err != nil {
		//removes images stored since transaction failed
		saver.CleanUp(paths)
		transaction.Rollback()
		return err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return nil

}
