package main

import (
	"net/http"
)

//Routes setup all routes in one file
func (s *Server) Routes() {

	//file Server
	fileServer := http.FileServer(http.Dir("./cmd/"))
	s.Router.PathPrefix("/static/").Handler(fileServer)

	//Products
	s.Router.HandleFunc("/products/{id}", s.HandleProductDelete).Methods("DELETE")
	s.Router.HandleFunc("/products/{id}", s.HandleProductGet).Methods("GET")
	s.Router.HandleFunc("/products", s.HandleProductCreate).Methods("POST")
	s.Router.HandleFunc("/search", s.HandleProductSearch).Methods("GET")

	// s.Router.HandleFunc("/categories", s.HandleProducts)
	// s.Router.HandleFunc("/collections", s.HandleProducts)
	// s.Router.HandleFunc("/users", s.HandleProducts)
	// s.Router.HandleFunc("/orders", s.HandleProducts)
}

//http.StripPrefix("/", fileServer)
