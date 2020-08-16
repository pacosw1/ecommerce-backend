package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Server unites all dependencies and prevents global varibles / singletons
type Server struct {
	Router *mux.Router
	Logger *Logger
	Models *Models
}

//NewServer initializes server without any dependencies except routes
func NewServer() *Server {
	serv := &Server{}
	serv.Router = mux.NewRouter()

	//file server for images

	//setup routes for serving
	serv.Routes()

	return serv
}

//ServeFiles serves static directory

//Listen serves to start listing for requests on selected Port
func (s *Server) Listen(port int) error {

	portAddr := ":" + strconv.Itoa(port)
	err := http.ListenAndServe(portAddr, s.Router)
	return err
}

//ServeHTTP allows server struct to become a http handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
