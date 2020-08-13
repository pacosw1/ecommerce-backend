package web

import (
	"net/http"
	"project-z/db"

	"github.com/gorilla/mux"
)

//Server unites all dependencies and prevents global varibles / singletons
type Server struct {
	Database *db.Database
	Router   *mux.Router
}

//NewServer initializes server without any dependencies except routes
func NewServer() *Server {
	serv := &Server{}
	serv.Router = mux.NewRouter()

	//setup routes for serving
	serv.Routes()

	return serv
}

//Listen serves to start listing for requests on selected Port
func (s *Server) Listen(port string) {

	http.ListenAndServe(":8080", s.Router)
}

//ServeHTTP allows server struct to become a http handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
