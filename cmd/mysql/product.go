package mysql

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	saver "project-z/cmd/image-saver"
	"project-z/cmd/models"
	"strings"
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

//Get gets a product
func (m *ProductModel) Get(id string) (*models.Product, error) {

	query := `SELECT * FROM products WHERE productID = ?`
	row := m.DB.QueryRow(query, id)
	var p models.Product

	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Stock, &p.Price, &p.SalePrice, &p.Created)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	query = "SELECT * FROM Images WHERE Images.productID = ?"
	rows, err := m.DB.Query(query, id)

	defer rows.Close()

	images := []*models.Image{}

	for rows.Next() {

		var img models.Image

		err = rows.Scan(&img.ImageID, &img.ProductID, &img.Path, &img.Thumbnail)
		if err != nil {
			return nil, err
		}

		url := strings.Split(img.Path, "/")
		path := strings.Join(url[1:], "/")

		img.Path = path

		images = append(images, &img)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	p.Images = images

	return &p, nil
}

//Search searches product by name
func (m *ProductModel) Search(name string) ([]*models.Product, error) {

	query := "SELECT * FROM products WHERE name LIKE ?"

	keyword := `%` + name + `%`

	rows, err := m.DB.Query(query, keyword)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := []*models.Product{}

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

	for _, p := range results {

		row := m.DB.QueryRow("SELECT * FROM images WHERE thumbnail = ? AND productId = ?", true, p.ID)

		var img models.Image

		err = row.Scan(&img.ImageID, &img.ProductID, &img.Path, &img.Thumbnail)

		if err != nil {
			return nil, err
		}

		url := strings.Split(img.Path, "/")
		path := strings.Join(url[1:], "/")

		img.Path = path

		p.Images = append(p.Images, &img)

	}

	return results, nil

}

//Remove removes product from webpage
func (m *ProductModel) Remove(id string) error {

	tx, err := m.DB.Begin()

	if err != nil {
		return err
	}

	query := "SELECT path FROM Images WHERE productID = ?"

	rows, err := tx.Query(query, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	defer rows.Close()

	paths := []string{}

	for rows.Next() {
		var path string

		err = rows.Scan(&path)

		if err != nil {
			tx.Rollback()
			return err
		}

		paths = append(paths, path)

	}

	query = "DELETE FROM IMAGES WHERE productID = ?"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "DELETE FROM Products WHERE productID = ?"
	_, err = tx.Exec(query, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	//if database successfully deleted then transfer images from disk to trash bin
	err = saver.DeleteImages(paths)
	if err != nil {
		//restore images to disk
		//TODO
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil

}
