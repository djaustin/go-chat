package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	// Allows a function to be run only once per instance of the struct
	once sync.Once
	// filename within templates dir to use as template
	filename string
	// template struct used to store parsed template
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// Parse the template for the filename provided. Only needs to be run once for each template
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// We don't have any substitutions in the template yet so pass no data
	t.templ.Execute(w, nil)
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
