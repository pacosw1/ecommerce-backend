package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"project-z/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql" // New import
)

func main() {

	//add command line flag to setup port dynamically
	addr := flag.Int("port", 4000, "HTTP network address")

	dsn := flag.String("dsn", "admin:o23pvfhs@/project_zilo?parseTime=true", "MySQL data source name")

	//make sure flags are valid
	flag.Parse()

	//Custom Loggers
	infoLog := log.New(os.Stdout, "INFO/t", log.Ldate|log.Ltime)
	//for erros
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//create a new empty server instance
	server := NewServer()

	//open database connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	//initialize models and give them access to database
	models := &Models{
		Product:  &mysql.ProductModel{DB: db},
		Category: &mysql.CategoryModel{DB: db},
	}

	//inject dependencies to server
	server.Models = models

	infoLog.Printf("Conntected to SQL Database")

	infoLog.Printf("Starting Web Server on PORT :%d", *addr)

	err = server.Listen(*addr)
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

//create a new DB instance and assign it to server
// database := db.NewDatabase()

// //attempt to Mongo Atlas Cluster
// cancel := server.Database.Connect(uri, "zilo")
// defer cancel()
// uri := "mongodb+srv://pacosw:o23pvfhs@cluster0.hvlbh.gcp.mongodb.net/zilo?retryWrites=true&w=majority"
