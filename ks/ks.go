package main

/*
TODO
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
	"strconv"
	"math/rand"
	"time"
	"github.com/gorilla/mux"
)

var db *sql.DB

var templates = template.Must(template.ParseFiles("index.html"))

var nonCryptoRand = rand.New( rand.NewSource( time.Now().UTC().UnixNano() ) )


func setupDatabase( hostName string, pgPassword string ) { // TODO pass in port, username
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


func getKey( keyID int64, userID int64 ) (string) {
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
		return keyVal;
	}
	return ""; 
}


func createKey(  userID int64, keyVal string ) (int64) {
	var keyID int64 = nonCryptoRand.Int63();
	var err error;

	var stmt [3]*sql.Stmt;
	var cmd  [3]string = [3]string{"INSERT INTO keys (kID, kVal, oID) VALUES ($1, $2::bytea,$3)",
		"INSERT INTO keyUsers (kID,uID) VALUES ($1,$2)",
		"INSERT INTO keyAdmins (kID,uID) VALUES ($1,$2);" }
	// note if using mySQL use ? but Postgres is $1 in prepare statements

	for i := range cmd {	
		stmt[i], err = db.Prepare( cmd[i] )
		if err != nil {
			log.Println("sql fatal error in createKey prep for", cmd[i])
			log.Fatal(err)
		}
	}

	for i := range cmd {
		if i == 0 {
			_,err = stmt[i].Exec(keyID,keyVal,userID)
		} else {
			_,err = stmt[i].Exec(keyID,userID)
		}
		
		if err != nil {
			log.Println("sql error in createKey",err)
			return 0;
		}
	}
	
	return keyID;
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
	vars := mux.Vars(r)

	var keyID int64 = 0;
	var userID int64 = 1;

	keyID,err = strconv.ParseInt(  vars["keyID"] , 0, 64 );
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Println("GET keyID=", keyID, " userID=",userID )

	io.WriteString(w, getKey(keyID,userID) )
}


func createKeyHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
	err := r.ParseForm()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }

	var keyVal string = r.FormValue("keyVal");
	var userID int64 = 1;

	var keyID int64 = createKey( userID, keyVal );

	log.Println("POST keyID=", keyID, "userID=",userID )

	io.WriteString(w, "{ \"keyID\": " + strconv.FormatInt(keyID,10) +" }" )
}


func main() {
	var err error

	// get all the configuration data 
	if len( os.Args ) != 2 {
		log.Fatal( "must pass database hostname on CLI" );	
	}
	var hostName string = os.Args[1]
	
	var pgPassword string = os.Getenv("SECM_DB_SECRET")
	if len( pgPassword ) < 1 {
		log.Fatal( "must set environ variable SECM_DB_SECRET" );	
	}

	// set up the DB 
	setupDatabase(hostName,pgPassword)
        defer db.Close()

	// Check DB is alive 
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// set up the routes 
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/v1/key/{keyID}", searchKeyHandler).Methods("GET")
	router.HandleFunc("/v1/key", createKeyHandler).Methods("POST")
	http.Handle("/", router)

	// run the web server 
        http.ListenAndServe(":8080", nil)
}
