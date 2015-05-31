package main

/*
TODO
*/

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("index.html"))


func mainHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	type PageData struct {
		Email string
		UserID int64
	}

	data := PageData{ Email: "no", UserID:0  }
	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func main() {
	// get all the configuration data
	if len(os.Args) != 2 {
		log.Fatal("must pass database hostname on CLI")
	}

		// set up the routes
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")

	http.Handle("/", router)

	// run the web server
	http.ListenAndServe(":8081", nil)
}
