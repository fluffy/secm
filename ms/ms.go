package main

/*
TODO

Connect up seqNo generation to DB 
Connect msg store to DB 
HTML test page doing the auth dance  
fetch integrity key from KS 
client JS to sign the msg 
verify signature of posted msg
webhooks to let notification server know about new msg
*/

import (
	"log"
	"net/http"
	"html/template"
	"os"
	"github.com/gorilla/mux"
	"strconv"
	"io"
	"io/ioutil"
)

var templates = template.Must(template.ParseFiles("index.html"))

var ksURL string = "https://localhost:443/"

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	type PageData struct {
		// Need to start with upper case letter 
		KSUrl string
	}

	data := PageData{ KSUrl: ksURL  }
	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		return
	}
	
	
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := string(contents)

	
	var seqNo int64 = 1
	
	log.Println("POST createMsg: keyID=", keyID, "seqNo=", seqNo, "msg=", msg)
	io.WriteString(w, "{ \"msgID\": "+ "\""+strconv.FormatInt(keyID, 10)+"-"+strconv.FormatInt(seqNo, 10)+ "\"" +" }")
}

func getMsgHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	log.Println("In getMsgHandler")
	
	var keyID int64
	keyID, err = strconv.ParseInt( vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var seqNo int64
	seqNo, err = strconv.ParseInt( vars["seqNo"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	log.Println("GET msg: msgID=", keyID,"-",seqNo)

	io.WriteString(w,   "{ \"msg\": \""  + "hello" + "\" " + "}" )
}

func main() {
	// get all the configuration data
	if len(os.Args) != 2 {
		log.Fatal("must pass KS URL on CLI")
	}

	ksURL = os.Args[1]
	
	// set up the routes
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/v1/msg/{keyID}", createMsgHandler).Methods("POST")
	router.HandleFunc("/v1/msg/{keyID}-{seqNo}", getMsgHandler).Methods("GET")

	http.Handle("/", router)

	// run the web server
	http.ListenAndServe(":8081", nil)
}
