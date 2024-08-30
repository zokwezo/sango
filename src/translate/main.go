package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var dbDvdRental *sql.DB

func main() {
	log.SetFlags(log.Lshortfile)
	var err error
	dbDvdRental, err = sql.Open(
		"postgres",
		"host=localhost dbname=dvdrental connect_timeout=5 statement_timeout=30 sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", onLoadHandler)
	http.HandleFunc("/fetch-actors/", fetchActorsHandler)

	fmt.Println("Navigate browser to http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Actor struct {
	FirstName string
	LastName  string
}
type Actors []Actor
type ActorsMap map[string]Actors

func onLoadHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("translate.html"))
	actorsMap := ActorsMap{"Actors": {}}
	err := tmpl.Execute(w, actorsMap)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchActorsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("translate.html"))
	time.Sleep(500 * time.Millisecond)
	firstNameToMatch := r.PostFormValue("firstName")
	lastNameToMatch := r.PostFormValue("lastName")
	actorsMap := GetActorsMapFromDatabase(dbDvdRental, firstNameToMatch, lastNameToMatch)
	err := tmpl.ExecuteTemplate(w, "actor-list-elements", actorsMap)
	if err != nil {
		log.Fatal(err)
	}
}

func GetActorsMapFromDatabase(db *sql.DB, firstNameToMatch, lastNameToMatch string) ActorsMap {
	rows, err := db.Query(makeQuerySql(firstNameToMatch, lastNameToMatch))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	actors := make(Actors, 0, 256)
	for rows.Next() {
		var actor Actor
		err := rows.Scan(&actor.FirstName, &actor.LastName)
		if err != nil {
			log.Fatal(err)
		}
		actors = append(actors, actor)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	actorsMap := make(ActorsMap)
	actorsMap["Actors"] = actors
	return actorsMap
}

func makeQuerySql(firstNameToMatch, lastNameToMatch string) string {
	makeWhereClause := func() string {
		whereClause := ""
		delim := " WHERE"
		if firstNameToMatch != "" {
			whereClause += delim +
				" REGEXP_LIKE(first_name, '" +
				strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(firstNameToMatch), "\""), "\"") +
				"')"
			delim = " AND"
		}
		if lastNameToMatch != "" {
			whereClause += delim +
				" REGEXP_LIKE(last_name, '" +
				strings.TrimSuffix(strings.TrimPrefix(strconv.Quote(lastNameToMatch), "\""), "\"") +
				"')"
			// delim = " AND"
		}
		return whereClause
	}

	querySql := "SELECT first_name, last_name FROM actor"
	if whereClause := makeWhereClause(); whereClause != "" {
		querySql += whereClause
	}
	return querySql
}
