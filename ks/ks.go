package main

/*
TODO
set up webserver
create new keys
add user for key 
add admin for key 
front end with appache
appache user auth 
dockerize 
api doc
figure out way to load secrets 
encrypt keys in DB 
DB backup

*/

import (
        "database/sql"
	_ "github.com/lib/pq"
        "log"
	"os"
	 "html/template"
	"net/http"
	"io"
	"strings"
	"strconv"
)

var db *sql.DB

var templates = template.Must(template.ParseFiles("index.html"))


func setupDatabase( hostName string, pgPassword string ) { // todo pass in hostname, port, username
        var err error

        // set up DB
        db, err = sql.Open("postgres", "password='"+pgPassword+"' user=postgres dbname=postgres host="+hostName+" port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// a postgess bigint is signed 64 bit int 
	sqlSetup := `
        CREATE TABLE IF NOT EXISTS keys ( kID BIGINT NOT NULL, kVal bytea NOT NULL ,  oID BIGINT NOT NULL, PRIMARY KEY( kID ) );
        CREATE TABLE IF NOT EXISTS keyUsers ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );
        CREATE TABLE IF NOT EXISTS keyAdmins ( kID BIGINT NOT NULL, uID BIGINT NOT NULL , PRIMARY KEY( kID,uID ) );
        `
        _, err = db.Exec(sqlSetup)
        if err != nil {
                log.Println("sql fatal error in setupDatabase")
                log.Printf("%q\n", err)
        }

	log.Printf("Setup Database\n")
}


func getKey( keyID, userID int64 ) (string) {
	// note if using mySQL use ? but Postgres is $1 in prepare statements 
	stmt, err := db.Prepare(
		"SELECT keys.kVal  FROM keyUsers JOIN keys ON  keys.kID = keyUsers.kID WHERE keyUsers.uID = $2 AND keyUsers.kID = $1")
	if err != nil {
		log.Println("sql fatal error in getKey prep")
		log.Fatal(err)
	}
	var keyVal string
	err = stmt.QueryRow(keyID,userID).Scan( &keyVal ) 
	switch {
	case err == sql.ErrNoRows:
            log.Printf("no key found")
	case err != nil:
		log.Println("sql fatal error in getKey querry")
		log.Fatal(err)
	default:
		log.Println("got key " + keyVal)
		return keyVal;
	}
	return ""; 
}


func mainHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }

        type PageData struct {
                Junk string
        }
        data := PageData{Junk: "nothing"}
        err := templates.ExecuteTemplate(w, "index.html", data)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}


func searchKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	
        w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("got URL path " + r.URL.Path )

	elements := strings.Split( r.URL.Path  , "/" );
	// note since the URL starts with a /, we get an empty element as first thing in array
	
	if len(elements) != 4  {
		http.NotFound(w, r)
                return
	}
	
	var keyID int64 = 0;
	var userID int64 = 1;

	keyID,err = strconv.ParseInt( elements[3] , 0, 64 );
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Println("GET keyID=", keyID, " userID=",userID )
	
	io.WriteString(w, getKey(keyID,userID) )
}


func main() {
	var err error
	
	var pgPassword string = os.Getenv("SECM_DB_SECRET")
	var hostName string = os.Args[1]
	
	setupDatabase(hostName,pgPassword)
        defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/v1/key/", searchKeyHandler)

        http.ListenAndServe(":8080", nil)
}
