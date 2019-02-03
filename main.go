package main

import (
	"github.com/gorilla/mux"
	"github.com/jbuchbinder/gopnm"
	"image"
	"log"
	"net/http"
	"sync"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("POST").
		Path("/convert").
		HandlerFunc(convertImage)

	err := http.ListenAndServe(":8888", router)
	if err != nil {
		log.Println("Server HTTP: ", err.Error())
	}
}

func convertImage(w http.ResponseWriter, r *http.Request) {
	log.Println("[", r.Host, "] Incoming request.")

	file, _, err := r.FormFile("image")

	if err != nil {
		log.Println("[", r.Host, "] Param error: ", err.Error())
		http.Error(w, "Param \"image\" error: "+err.Error(), 400)

		return
	}

	img, err := pnm.Decode(file)

	if err != nil {
		log.Println("[", r.Host, "] File decode error: ", err.Error())
		http.Error(w, "File decode error: "+err.Error(), 400)

		return
	}

	newImage := async(&img)

	w.Header().Set("Content-Disposition", "attachment; filename=converted.pgm")
	w.Header().Set("Content-Type", "image/x-portable-graymap")

	err = pnm.Encode(w, newImage, pnm.PGM)

	if err != nil {
		log.Println("[", r.Host, "] File encode error: ", err.Error())
		http.Error(w, "File encode error: "+err.Error(), 400)

		return
	}

	log.Println("[", r.Host, "] Response sent.")
}

func async(oldImage *image.Image) image.Image {

	newImage := image.NewGray(image.Rectangle{Max: image.Point{X: 1024, Y: 1024}})

	countThreads := 1024
	var wg sync.WaitGroup
	wg.Add(countThreads)

	linesPerThread := (int)(1024 / countThreads)
	task := make([]Task, countThreads)

	for i := 0; i < countThreads; i++ {
		task[i] = Task{&wg, linesPerThread, i * linesPerThread, oldImage, newImage}
		task[i].BeginConvolution()
	}

	wg.Wait()

	return newImage
}
