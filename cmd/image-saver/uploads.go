package saver

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
)

//WriteImg stores image to disk
func WriteImg(path string, img image.Image) (bool, error) {

	//create new empty file to store image
	f, err := os.Create(path)

	if err != nil {
		fmt.Printf("Could not create file in os")
		return false, err
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
		return false, err
	}

	return true, nil
}

//SaveImagesToDisk takes in http file header, optimizes and saves them to disk
func SaveImagesToDisk(directory string, files []*multipart.FileHeader, thumbnail string) (string, []string, error) {

	var paths []string

	primaryPath := ""

	for _, header := range files {

		//save image to given dir and return its saved path
		display, path, err := SaveImage(header, directory, thumbnail)

		if err != nil {
			return "", nil, err
		}

		if display {
			primaryPath = path
		}

		paths = append(paths, path)
	}

	return primaryPath, paths, nil

}

//ValidateImageFiles validates right image format
func ValidateImageFiles(files []*multipart.FileHeader) error {

	if len(files) < 1 || len(files) > 10 {
		return errors.New("Invalid request, 0 images or more than 10")
	}
	for _, header := range files {

		file, err := header.Open()

		if err != nil {
			return err
		}

		defer file.Close()

		//validate file type
		err = ValidateImage(file)

		if err != nil {
			return err
		}
	}

	return nil
}

//CleanUp deletes images from disk if transaction failed
func CleanUp(paths []string) error {

	for _, path := range paths {

		err := os.Remove(path)

		if err != nil {
			return err
		}
	}

	return nil
}

//SaveImage saves image into specified directory
func SaveImage(header *multipart.FileHeader, directory string, thumbnail string) (bool, string, error) {

	isThumbNail := thumbnail == header.Filename

	//open file
	file, err := header.Open()

	if err != nil {
		return false, "", err
	}

	defer file.Close()

	// //validate file type
	// err = validateImage(file)

	// if err != nil {
	// 	return "", err
	// }

	// //rewind pointer since we read the heater to prove validity
	// file.Seek(0, 0)

	//convert file to image
	img, _, err := DecodeImage(file)

	if err != nil {
		return false, "", err
	}

	//resize the image for optimal web size
	img = resize.Resize(1500, 0, img, resize.Lanczos3)

	//generate a random image identifier name
	uid := uuid.NewV4()
	path := fmt.Sprintf("%s/%s.%s", directory, uid, "jpg")

	//write optimized image to storage
	_, err = WriteImg(path, img)

	if err != nil {
		return false, "", err
	}

	return isThumbNail, path, nil
}

//DecodeImage decodes image
func DecodeImage(file multipart.File) (image.Image, string, error) {

	img, format, err := image.Decode(file)

	if err != nil {
		return nil, "", err
	}

	return img, format, err
}

//ValidateImage validates the image content type to check image integrity and validity
func ValidateImage(img multipart.File) error {

	const (
		jpg = "image/jpeg"
		png = "image/png"
	)

	buffer := make([]byte, 512)

	_, err := img.Read(buffer)

	if err != nil {
		return err
	}

	contentType := http.DetectContentType(buffer)

	switch contentType {
	case jpg:
		return nil
	case png:
		return nil
	default:
		return errors.New("Unsupported Image Format " + contentType + ". Only JPG and PNG supported")
	}

}
