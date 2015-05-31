package main

/*
TODO

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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var templates = template.Must(template.ParseFiles("index.html"))

var ksURL string = "https://localhost:443/"

var session *mgo.Session
var msgCollection *mgo.Collection
var countCollection *mgo.Collection

type Message struct {
	Id string
	Data string
}

type Count struct {
	KeyID int64
	SeqNum int64
}


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

	var seqNum int64 = getNextSeq(keyID)
	var msgID string = strconv.FormatInt(keyID, 10) + "-" + strconv.FormatInt(seqNum, 10) 
	
	err = msgCollection.Insert( &Message{ msgID,msg } )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Println("POST createMsg: keyID=", keyID, "seqNum=", seqNum )
	
	io.WriteString(w, "{ \"seqNum\": "+ "\""+strconv.FormatInt(seqNum, 10)+ "\"" +" }")
}

func getMsgHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	var keyID int64
	keyID, err = strconv.ParseInt( vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var seqNum int64
	seqNum, err = strconv.ParseInt( vars["seqNum"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var msgID string = strconv.FormatInt(keyID, 10) + "-" + strconv.FormatInt(seqNum, 10) 
	
	log.Println("GET msg: msgID=", msgID )
	result := Message{}
	err = msgCollection.Find( bson.M{ "id": msgID } ).One(&result)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		http.NotFound(w, r)
		return
	}	

	io.WriteString(w,  result.Data  )
}

func getNextSeq(keyID int64) (int64) {
	var err error
	doc := Count{}
	
	change := mgo.Change{
        Update: bson.M{"$inc": bson.M{"seqnum": 1}},
        ReturnNew: true,
	}
	_ , err = countCollection.Find( bson.M{"keyid": keyID} ).Apply(change, &doc)
	if err != nil {
		// key did not exist, create a new one
		err = countCollection.Insert( &Count{ keyID,1 } )
		if err != nil {
			log.Fatal( "Could not create count" , err );
			return 0
		}
		return 1
	}

	return doc.SeqNum;
}

func main() {
	var err error
	
	// get all the configuration data
	if len(os.Args) != 3 {
		log.Fatal("CLI must have keyServerUrl then mongoHostName ")
	}

	ksURL = os.Args[1]
	var mongoHost string = os.Args[2]

	session, err = mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	msgCollection = session.DB("secm").C("messages")
	countCollection = session.DB("secm").C("counts")

	var keyID int64 = 12345
	log.Println( "next seq=" , getNextSeq(keyID) )

	
	// set up the routes
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/v1/ch/{keyID}", createMsgHandler).Methods("POST")
	router.HandleFunc("/v1/msg/{keyID}-{seqNum}", getMsgHandler).Methods("GET")

	http.Handle("/", router)

	// run the web server
	log.Println("Ready for requests")
	http.ListenAndServe(":8081", nil)
}
