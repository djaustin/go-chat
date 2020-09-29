package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/djaustin/go-chat/trace"
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
	t.templ.Execute(w, r)
}

func main() {
	addr := flag.String("addr", ":8080", "The address of the application")
	flag.Parse()

	mainRoom := newRoom()
	mainRoom.tracer = trace.New(os.Stdout)
	go mainRoom.run()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", mainRoom)

	log.Println("Starting application on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
