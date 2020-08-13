package main

import (
	"project-z/db"
	"project-z/web"
)

func main() {

	uri := "mongodb+srv://pacosw:o23pvfhs@cluster0.hvlbh.gcp.mongodb.net/zilo?retryWrites=true&w=majority"

	//create a new empty server instance
	server := web.NewServer()

	//create a new DB instance and assign it to server
	database := db.NewDatabase()
	server.Database = database

	//attempt to Mongo Atlas Cluster
	cancel := server.Database.Connect(uri, "zilo")

	defer cancel()

	server.Listen("8080")

}
