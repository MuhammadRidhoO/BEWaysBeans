package middleware

import (
	"BEWaysBeans/dto"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(next http.HandlerFunc, formInputName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := r.FormFile(formInputName)

		if err != nil {
			// fmt.Println(err)
			// add filename to ctx
			ctx := context.WithValue(r.Context(), "Error", "Image Failed to Upload?")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		defer file.Close()
		// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		// fmt.Printf("File Size: %+v\n", handler.Size)
		// fmt.Printf("MIME Header: %+v\n", handler.Header)
		const MAX_UPLOAD_SIZE = 10 << 20 // 10MB
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Code: http.StatusBadRequest, Message: "Max size in 1mb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			fmt.Println(err, "path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filename := data[8:] // split uploads/

		// add filename to ctx
		ctx := context.WithValue(r.Context(), "dataFile", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func UploadTripImage(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handling dan parsing data dari form data yang ada data file nya. Argumen 1024 pada method tersebut adalah maxMemory sebesar 1024byte, apabila file yang diupload lebih besar maka akan disimpan di file sementara
		if err := r.ParseMultipartForm(1024); err != nil {
			panic(err.Error())
		}

		var arrImages []string

		files := r.MultipartForm.File["images"]
		for _, f := range files {
			// mengambil file dari form
			file, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			defer file.Close()

			// Apabila format file bukan .jpg, .jpeg atau .png, maka tampilkan error
			if filepath.Ext(f.Filename) != ".jpg" && filepath.Ext(f.Filename) != ".jpeg" && filepath.Ext(f.Filename) != ".png" {
				w.WriteHeader(http.StatusBadRequest)
				response := dto.ErrorResult{
					Code:    http.StatusBadRequest,
					Message: "The provided file format is not allowed. Please upload a JPG, JPEG or PNG image",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// create empty context
			var ctx = context.Background()

			// setup cloudinary credentials
			var CLOUD_NAME = os.Getenv("CLOUD_NAME")
			var API_KEY = os.Getenv("API_KEY")
			var API_SECRET = os.Getenv("API_SECRET")

			// create new instance of cloudinary object using cloudinary credentials
			cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

			// Upload file to Cloudinary
			resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "DeweTour"})
			if err != nil {
				fmt.Println(err.Error())
			}
			// cek respon dari cloudinary
			// fmt.Println("respon from cloudinary", resp)

			// arrImages = append(arrImages, fileLocation)
			arrImages = append(arrImages, resp.SecureURL)
		}

		// membuat sebuah context baru dengan menyisipkan value di dalamnya, valuenya adalah array sring yang berisikan url img yang didapat dari cloudinary
		ctx := context.WithValue(r.Context(), "arrImages", arrImages)

		// mengirim nilai context ke object http.HandlerFunc yang menjadi parameter saat fungsi middleware ini dipanggil
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
