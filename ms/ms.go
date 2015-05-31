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
	"strconv"
	"io"
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

func createMsgHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var keyID int64
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	keyID, err = strconv.ParseInt( vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var msg string = r.FormValue("msg")

	var seqNo int64 = 1
	
	log.Println("POST createMsg: keyID=", keyID, "seqNo=", seqNo, "msg=", msg)
	io.WriteString(w, "{ \"msgID\": "+ "\""+strconv.FormatInt(keyID, 10)+"-"+strconv.FormatInt(seqNo, 10)+ "\"" +" }")
}

func main() {
	// get all the configuration data
	if len(os.Args) != 2 {
		log.Fatal("must pass database hostname on CLI")
	}

		// set up the routes
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/v1/msg/{keyID}", createMsgHandler).Methods("POST")
	//router.HandleFunc("/v1/msg/{msgID}", getMsgHandler).Methods("GET")

	http.Handle("/", router)

	// run the web server
	http.ListenAndServe(":8081", nil)
}
