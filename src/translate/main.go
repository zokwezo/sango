package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("translate.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		err := tmpl.Execute(w, films)
		if err != nil {
			log.Fatal(err)
		}
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		tmpl := template.Must(template.ParseFiles("translate.html"))
		err := tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film/", h2)

	fmt.Println("Navigate browser to http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
