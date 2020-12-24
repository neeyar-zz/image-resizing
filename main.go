package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"log"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"image/gif"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return 
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
 
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}



func transformFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return 
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)

	filetype := handler.Header["Content-Type"][0]
    filename := handler.Filename
    filename_len := len(filename)
    filename_without_type := filename[:filename_len-4]

    switch filetype {
    case "image/jpeg", "image/jpg":
            fmt.Println("File Type: ", filetype)

            img, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			m := resize.Resize(1280, 0, img, resize.Lanczos3)

			out, err := os.Create("resized-" + filename_without_type + ".jpg")
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			jpeg.Encode(out, m, nil)

			fmt.Fprintf(w, "The File has been successfully resized and saved\n")

    case "image/gif":
            fmt.Println("File Type: ", filetype)
            fmt.Println("the animation resizing aren't supported")

            img, err := gif.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()


			m := resize.Resize(1280, 0, img, resize.Lanczos3)

			out, err := os.Create("resized-" + filename_without_type + ".gif")
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			gif.Encode(out, m, nil)	

			fmt.Fprintf(w, "The File has been successfully resized and saved\n")

         
    case "image/png":
            fmt.Println("File Type: ", filetype)

            img, err := png.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			m := resize.Resize(1280, 0, img, resize.Lanczos3)

			out, err := os.Create("resized-" + filename_without_type + ".png")
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			png.Encode(out, m)	

			fmt.Fprintf(w, "The File has been successfully resized and saved\n")

    default:
            fmt.Println("unknown file type uploaded")
            fmt.Fprintf(w, "unknown file type uploaded. supported types: *jpeg, *png, *gif\n")
    }


}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/transform", transformFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("The server is running")
	setupRoutes()
}

