package web

import (
	"context"
	"encoding/json"
	"net/http"
	"project-z/db"
	"time"
)

//HandleProductCreate serves to create a new product and store it in database
func (s *Server) HandleProductCreate(w http.ResponseWriter, r *http.Request) {
	//set headers to expect json presonse from server
	w.Header().Set("Content-Type", "application/json")

	var product db.Product

	//decode from json to Product struct
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&product)
	if err != nil {
		panic(err)
	}

	//get reference to db
	db := s.Database.Root

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// insert document into mongo Collection
	res, err := db.Collection("products").InsertOne(ctx, product)

	// check for errors
	if err != nil {
		panic(err)
	}

	//send created product as response

	json.NewEncoder(w).Encode(res)

}

// msg := struct {
// 	Message string `json:"message"`
// }{
// 	Message: "Created Successfully",
// }

//category handlers

//
