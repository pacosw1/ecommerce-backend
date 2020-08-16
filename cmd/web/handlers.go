package main

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	"net/http"
	saver "project-z/cmd/image-saver"
	"project-z/cmd/models"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
)

//HandleProductCreate serves to create a new product and store it in database
func (s *Server) HandleProductCreate(w http.ResponseWriter, r *http.Request) {

	logger := s.Logger
	decoder := schema.NewDecoder()
	validate := validator.New()

	//set headers to expect json presonse from server
	w.Header().Set("Content-Type", "application/json")

	//if request is invalid
	if r.Method != http.MethodPost {
		//show caller allowed calls to this endpoint
		w.Header().Set("Allow", http.MethodPost)
		logger.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//max file bundle upload size (100 MB)
	var maxMb int64 = 100

	maxBytes := maxMb * 1024

	//parse form data that includes files
	err := r.ParseMultipartForm(maxBytes)

	if err != nil {
		http.Error(w, "Total Image Upload Exceeded, more than 100 MB", 400)
		return
	}

	//abstract images from the request
	images := r.MultipartForm.File["images"]

	//validate images
	err = saver.ValidateImageFiles(images)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var product models.Product
	//abtract regular json key value pairs
	decoder.Decode(&product, r.MultipartForm.Value)

	thumbnail := r.FormValue("thumbnail")

	if thumbnail == "" {
		http.Error(w, "Invalid thumbnail", 400)
		return
	}

	exists := false
	for _, header := range images {
		if header.Filename == thumbnail {
			fmt.Print(thumbnail)
			exists = true
		}
	}

	if !exists {
		http.Error(w, "Make sure selected thumbnail is attached", 400)
		return
	}

	err = validate.Struct(product)

	//validation errors
	if err != nil {

		type ErrorResponse struct {
			Errors []map[string]string
		}

		resp := &ErrorResponse{}

		errors := make([]map[string]string, 0)

		for _, err := range err.(validator.ValidationErrors) {
			newErr := make(map[string]string)

			newErr["code"] = "Validation"
			newErr["message"] = err.Field() + " is " + err.Tag()

			errors = append(errors, newErr)
		}
		resp.Errors = errors

		http.Error(w, "Bad Request", http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return

	}

	//send validated data to our DB transaction function
	err = s.Models.Product.Insert(product, images, thumbnail)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode("Save Successfull")
	return
}

//HandleProductSearch search for a product based on its name
func (s *Server) HandleProductSearch(w http.ResponseWriter, r *http.Request) {
	//set headers to expect json presonse from server
	w.Header().Set("Content-Type", "application/json")

	//if request is invalid
	if r.Method != http.MethodGet {
		//show caller allowed calls to this endpoint
		w.Header().Set("Allow", http.MethodGet)
		return
	}

	r.ParseForm()

	keyword := r.FormValue("keyword")

	if keyword == "" {
		http.Error(w, "Invalid Keyword", 400)
		return
	}

	results, err := s.Models.Product.Search(keyword)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(results)
}
