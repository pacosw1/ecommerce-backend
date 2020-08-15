package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"project-z/pkg/models"
)

//HandleProductCreate serves to create a new product and store it in database
func (s *Server) HandleProductCreate(w http.ResponseWriter, r *http.Request) {

	logger := s.Logger
	//set headers to expect json presonse from server
	w.Header().Set("Content-Type", "application/json")

	//if request is invalid
	if r.Method != http.MethodPost {

		//show caller allowed calls to this endpoint
		w.Header().Set("Allow", http.MethodPost)
		logger.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var product models.Product

	//setup decoder
	decoder := json.NewDecoder(r.Body)

	//decode into product struct and check for erros
	err := decoder.Decode(&product)
	if err != nil {
		logger.serverError(w, err)
	}

	//access DB model
	productModel := s.Models.Product

	//Insert into database and check for err
	id, err := productModel.Insert(product)
	if err != nil {
		logger.serverError(w, err)
	}

	//if no error resulted, send new Record ID
	json.NewEncoder(w).Encode(id)

}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

//HandleImageCreate serves to store an image in database
func (s *Server) HandleImageCreate(w http.ResponseWriter, r *http.Request) {

	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	// image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	r.ParseMultipartForm(100 * 1024)

	type ImageUpload struct {
		Images []byte `json:"images"`
		Name   string `json:"name"`
	}

	const (
		jpg = "image/jpeg"
		png = "image/png"
	)

	var maxMb int64 = 10

	maxBytes := maxMb * 1024
	r.ParseMultipartForm(maxBytes)

	images := r.MultipartForm.File["images"]

	fmt.Print(len(images))
	if len(images) == 0 || len(images) > 10 {
		fmt.Printf("Invalid Image Count")
	}

	for _, header := range images {

		file, err := header.Open()

		if err != nil {
			fmt.Printf("File could not be opened")
			panic(err)
		}

		defer file.Close()

		_, valid, err := validateImage(file)

		if err != nil || !valid {
			fmt.Printf("File type not supported")
		}

		//rewind pointer since we read the heater to prove validity
		file.Seek(0, 0)

		//convert file to image
		img, err := jpeg.Decode(file)

		if err != nil {
			panic(err)
		}

		// fileName := "testfile"
		// format := "jpg"

		// path := fmt.Sprintf("/%s.%s", fileName, format)

		//create new empty file to store image
		f, err := os.Create("cmd/storage/test.jpg")

		if err != nil {
			fmt.Printf("Could not create file in os")
			panic(err)
		}

		defer f.Close()

		//compress file jpg only
		compression := jpeg.Options{
			Quality: 65,
		}

		//encode newly compressed image
		err = jpeg.Encode(f, img, &compression)

		if err != nil {
			fmt.Printf("Image encoding Error")
			panic(err)
		}

		json.NewEncoder(w).Encode(image.DecodeConfig)

	}

}

func encodePng(img *image.Image, fileName string) {

}

func encodeJpg(img *image.Image, fileName string, quality int) {

}

//validates the image content type to check image integrity and validity
func validateImage(img multipart.File) (string, bool, error) {

	buffer := make([]byte, 512)

	_, err := img.Read(buffer)

	if err != nil {
		return "", false, err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, true, nil

}
