package main

//Routes setup all routes in one file
func (s *Server) Routes() {
	s.Router.HandleFunc("/products", s.HandleProductCreate).Methods("POST")
	// s.Router.HandleFunc("/categories", s.HandleProducts)
	// s.Router.HandleFunc("/collections", s.HandleProducts)
	// s.Router.HandleFunc("/users", s.HandleProducts)
	// s.Router.HandleFunc("/orders", s.HandleProducts)
}
