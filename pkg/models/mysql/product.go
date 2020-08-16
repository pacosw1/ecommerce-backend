package mysql

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	saver "project-z/pkg/image-saver"
	"project-z/pkg/models"
)

//ProductModel type which wraps a sql.DB connection pool.
type ProductModel struct {
	DB *sql.DB
}

//Insert serves to add a product to the database
func (m *ProductModel) Insert(p models.Product, images []*multipart.FileHeader, thumbnail string) error {

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
	primaryPath, paths, err := saver.SaveImagesToDisk("cmd/static/images", images, thumbnail)

	if err != nil {
		transaction.Rollback()
		//removes images stored since transaction failed
		saver.CleanUp(paths)
		return err
	}

	var values []interface{}
	imageQuery := `INSERT INTO images (thumbnail, productId, path) VALUES `

	for index, path := range paths {

		row := "(?, ?, ?)"
		if index < len(paths)-1 {
			row += ", "
		}
		imageQuery += row
		values = append(values, path == primaryPath, productID, path)
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

//Search searches product by name
func (m *ProductModel) Search(name string) ([]*models.Product, error) {

	query := "SELECT * FROM products WHERE name LIKE ?"

	keyword := `%` + name + `%`

	rows, err := m.DB.Query(query, keyword)

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	defer rows.Close()

	var results []*models.Product

	//iterate thru results
	for rows.Next() {

		p := models.Product{}

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SalePrice, &p.Stock, &p.Created)

		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

		results = append(results, &p)

	}

	if err = rows.Err(); err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	type ProductSearch struct {
		Name      string
		ImagePath string
	}

	// searchRes := []*ProductSearch{}

	fmt.Print(results)

	for _, product := range results {
		res := ProductSearch{}
		res.Name = product.Name
		row := m.DB.QueryRow("SELECT path FROM images WHERE thumbnail = true")

		var path string
		row.Scan(&path)

		fmt.Print(row)

	}

	return results, nil

}
