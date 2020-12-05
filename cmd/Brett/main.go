package main

import (
	"log"
	"net/http"
	"time"

	"github.com/zischknall/Brett/pkg/storage"

	"github.com/gorilla/mux"
)

func main() {
	store, err := storage.GetFileStore("/tmp/media")
	if err != nil {
		log.Fatal("Unable to get FileStore", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(rootHTML))
	})

	r.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {
		uploadFile(w, r, store)
	}).Methods("POST")

	r.HandleFunc("/storage/{filename}", func(w http.ResponseWriter, r *http.Request) {
		retrieveFile(w, r, store)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func uploadFile(w http.ResponseWriter, r *http.Request, s storage.Store) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Failed passing file from form")

		return
	}

	hash, err := s.SaveFile(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to save file to disk")

		return
	}

	err = file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to close file")

		return
	}

	_, _ = w.Write([]byte(hash))
}

func retrieveFile(w http.ResponseWriter, r *http.Request, s storage.Store) {
	vars := mux.Vars(r)
	file, err := s.GetFileWithHash(vars["filename"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	http.ServeContent(w, r, "", time.Now(), file)
}

const rootHTML = `
<!DOCTYPE html>
<html>
<head>
<title>Brett</title>
</head>
<body>
<form method="post" action="/storage" enctype="multipart/form-data">
  <div>
    <label for="file">Choose a file</label>
    <input type="file" id="file" name="file" multiple>
  </div>
  <div>
    <button>Send the file</button>
  </div>
</form>
</body>
</html>
`
