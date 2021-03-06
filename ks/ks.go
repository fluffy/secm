package main

/*
TODO

- Need to sort out auth for iKey so that only user of key or cloud can read it 

move the returned ID to be base32 or base64 encoded 

move all prepares statements to DB setup time
defer stmt.Close() for all statements
encrypt keys in DB
DB backup
*/

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"crypto/sha1"
	"encoding/binary"
)

var db *sql.DB

var templates = template.Must(template.ParseFiles("index.html"))

var nonCryptoRand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func setupDatabase(hostName string, pgPassword string) { // TODO pass in port, username
	var err error

	// set up DB
	db, err = sql.Open("postgres", "password='"+pgPassword+"' user=postgres dbname=postgres host="+hostName+" port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// a postgess bigint is signed 64 bit int
	sqlSetup := `
        CREATE TABLE IF NOT EXISTS keys ( kID BIGINT NOT NULL, kVal bytea NOT NULL, ikVal bytea NOT NULL, oID BIGINT NOT NULL, PRIMARY KEY( kID ) );
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

func getKey(keyID int64, userID int64) string {
	// note if using mySQL use ? but Postgres is $1 in prepare statements
	stmt, err := db.Prepare(
		"SELECT keys.kVal  FROM keyUsers JOIN keys ON  keys.kID = keyUsers.kID WHERE keyUsers.uID = $2 AND keyUsers.kID = $1")
	if err != nil {
		log.Println("sql fatal error in getKey prep")
		log.Fatal(err)
	}
	var keyVal string
	err = stmt.QueryRow(keyID, userID).Scan(&keyVal)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no key found")
	case err != nil:
		log.Println("sql fatal error in getKey querry")
		log.Fatal(err)
	default:
		return keyVal
	}
	return ""
}

func getIKey(keyID int64 ) string {
	// note if using mySQL use ? but Postgres is $1 in prepare statements
	stmt, err := db.Prepare(
		"SELECT keys.ikVal FROM keys WHERE kID = $1")
	if err != nil {
		log.Println("sql fatal error in getIKey prep")
		log.Fatal(err)
	}
	var iKeyVal string
	err = stmt.QueryRow(keyID).Scan(&iKeyVal)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no key found")
	case err != nil:
		log.Println("sql fatal error in getKey querry")
		log.Fatal(err)
	default:
		return iKeyVal
	}
	return ""
}

func createKey(userID int64, keyVal string, iKeyVal string) int64 {
	
	//var max int64 =  1<<60; var min int64 = 1<<55; // for base 32 encoding 
	var max int64 =  922337203685477585; var min int64 = 100000000000000000; // for base 10 encoding
	var keyID int64 = min + nonCryptoRand.Int63n( max-min )
	
	var err error

	var stmt [3]*sql.Stmt
	var cmd [3]string = [3]string{"INSERT INTO keys (kID, kVal, ikVal, oID) VALUES ($1, $2::bytea,$3::bytea,$4)",
		"INSERT INTO keyUsers (kID,uID) VALUES ($1,$2)",
		"INSERT INTO keyAdmins (kID,uID) VALUES ($1,$2);"}
	// note if using mySQL use ? but Postgres is $1 in prepare statements

	for i := range cmd {
		stmt[i], err = db.Prepare(cmd[i])
		if err != nil {
			log.Println("sql fatal error in createKey prep for", cmd[i])
			log.Fatal(err)
		}
	}

	for i := range cmd {
		if i == 0 {
			_, err = stmt[i].Exec(keyID, keyVal, iKeyVal, userID)
		} else {
			_, err = stmt[i].Exec(keyID, userID)
		}

		if err != nil {
			log.Println("sql error in createKey", err)
			return 0
		}
	}

	return keyID
}

func addRole(keyID int64, userID int64, role string, roleID int64) error {
	var err error
	var stmt *sql.Stmt
	var cmd string

	switch {
	case role == "user":
		cmd = "INSERT INTO keyUsers (kID,uID) SELECT kID,$3 FROM keyAdmins WHERE keyAdmins.kID = $1 AND keyAdmins.uID = $2"
	case role == "admin":
		cmd = "INSERT INTO keyAdmins (kID,uID) SELECT kID,$3 FROM keys WHERE keys.kID = $1 AND keys.oID = $2"
	default:
		return errors.New("bad role")
	}

	stmt, err = db.Prepare(cmd)
	if err != nil {
		log.Println("sql fatal error in addRole prep for", cmd)
		log.Fatal(err)
	}

	_, err = stmt.Exec(keyID, userID, roleID)
	if err != nil {
		log.Println("sql error in addRole", err)
		return err
	}

	return nil
}

func getMeta(keyID int64, userID int64, meta string) ([]int64, error) {

	var cmd string

	switch {
	case meta == "owner":
		cmd = "SELECT keys.oID FROM keys JOIN  keyUsers ON keys.kID = keyUsers.kID WHERE keys.kID = $1 AND keyUsers.uID = $2"
	case meta == "admins":
		cmd = "SELECT keyAdmins.uID FROM keyAdmins JOIN keyUsers ON keyAdmins.kID = keyUsers.kID WHERE keyUsers.kID = $1 AND keyUsers.uID = $2"
	case meta == "users":
		cmd = "SELECT users.uID FROM keyUsers AS users JOIN keyUsers AS perm ON users.kID = perm.kID WHERE perm.kID = $1 AND perm.uID = $2"
	default:
		return []int64{}, errors.New("bad meta")
	}

	// note if using mySQL use ? but Postgres is $1 in prepare statements
	stmt, err := db.Prepare(cmd)
	if err != nil {
		log.Println("sql fatal error in getMeta prep meta=", meta)
		log.Fatal(err)
	}

	rows, err := stmt.Query(keyID, userID)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("no data found") // key probably does not exist for this to happen
	case err != nil:
		log.Println("sql fatal error in getMeta querry")
		log.Fatal(err)
	}
	defer rows.Close()

	var ret []int64 = []int64{}

	for rows.Next() {
		var uID int64
		err := rows.Scan(&uID)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, uID)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ret, nil
}


func getUserID(  r *http.Request ) (int64,error) {

	var email string = r.Header.Get("OIDC_CLAIM_email")
	if email == "" {
		if true { // TODO remove - test mode
			return 1,nil
		}
		return 0, errors.New("bad authentication")
	}
	var verified string = r.Header.Get("OIDC_CLAIM_email_verified")
	if verified != "1" {
		return 0, errors.New("bad authentication")
	}
	
	data := []byte(email)
	hash := sha1.Sum(data)
	id := binary.BigEndian.Uint64( hash[0:8] )
	if id > 9223372036854775807 {
		id = id-9223372036854775807
	}
	v := int64( id )
	
	return v,nil
}


func mainHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var userID int64
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var email string = r.Header.Get("Oidc_claim_email")
	
	type PageData struct {
		Email string
		UserID int64
	}

	userID,err = getUserID(r)
	if err != nil {
		userID = 0
		email = "no-user"
	}
	
	data := PageData{ Email: email, UserID:userID  }
	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var userID int64 = 1
	userID,err = getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	log.Println("GET me: userID=", userID)

	io.WriteString(w,   "{ \"userID\":"  + " \"" + strconv.FormatInt( userID, 10) + "\" " + "}" )
}

func fetchKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	var keyID int64 = 0
	var userID int64 = 1

	userID,err = getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	keyID, err = strconv.ParseInt(vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Println("GET key: keyID=", keyID, " userID=", userID)

	io.WriteString(w, getKey(keyID, userID) )
}

func fetchIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)

	var keyID int64 = 0
	keyID, err = strconv.ParseInt( vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	log.Println("GET iKey: keyID=", keyID)

	io.WriteString(w, getIKey( keyID ) )
}

func createKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var keyVal string = r.FormValue("keyVal")
	var iKeyVal string = r.FormValue("iKeyVal")
	var userID int64 = 1

	userID,err = getUserID(r)
	if err != nil {
		log.Println("FAIL POST createKey due to unauthorized" )
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	var keyID int64 = createKey(userID, keyVal, iKeyVal)
	
	log.Println("POST createKey: keyID=", keyID, "userID=", userID)
	io.WriteString(w, "{ \"keyID\": "+ "\""+strconv.FormatInt(keyID, 10)+ "\"" +" }")
}

func addRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var keyID int64 = 0
	var userID int64 = 1
	var roleID int64 = 0

	userID,err = getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	keyID, err = strconv.ParseInt(vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	roleID, err = strconv.ParseInt(vars["roleID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var role string = vars["role"]

	switch {
	case role == "admin":
	case role == "user":
	default:
		http.NotFound(w, r)
		return
	}

	log.Println("POST addRole: keyID=", keyID, " userID=", userID, "roleID=", roleID, "role=", role)

	err = addRole(keyID, userID, role, roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}
}

func getKeyMetaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var err error

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var userID int64 = 1
	userID,err = getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	var keyID int64 = 0
	keyID, err = strconv.ParseInt(vars["keyID"], 0, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var meta string = vars["meta"]

	var vals []int64
	var idName string

	switch {
	case meta == "owner":
		idName = "ownerIDs"
	case meta == "admins":
		idName = "adminIDs"
	case meta == "users":
		idName = "userIDs"
	default:
		http.NotFound(w, r)
		return
	}

	log.Println("GET meta: keyID=", keyID, "userID=", userID, "meta=", meta)

	vals, err = getMeta(keyID, userID, meta)

	io.WriteString(w, "{ \""+idName+"\": [ ")
	for i := range vals {
		if i != 0 {
			io.WriteString(w, ",")
		}
		io.WriteString(w, "\"" + strconv.FormatInt(vals[i], 10) + "\"" )
	}
	io.WriteString(w, " ] }")
}

func main() {
	var err error

	// get all the configuration data
	if len(os.Args) != 2 {
		log.Fatal("must pass database hostname on CLI")
	}
	var hostName string = os.Args[1]
	var pgPassword string = os.Getenv("DB_ENV_POSTGRES_PASSWORD")
	if len(pgPassword) < 1 {
		log.Fatal("must set environ variable DB_ENV_POSTGRES_PASSWORD")
	}

	// set up the DB
	setupDatabase(hostName, pgPassword)
	defer db.Close()

	// Check DB is alive
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// set up the routes
	router := mux.NewRouter()
	router.HandleFunc("/", mainHandler).Methods("GET")
	router.HandleFunc("/v1/identity/me", meHandler).Methods("GET")
	router.HandleFunc("/v1/key", createKeyHandler).Methods("POST")
	router.HandleFunc("/v1/key/{keyID}", fetchKeyHandler).Methods("GET")
	router.HandleFunc("/v1/iKey/{keyID}", fetchIKeyHandler).Methods("GET")
	router.HandleFunc("/v1/key/{keyID}/{role}/{roleID}", addRoleHandler).Methods("PUT") // role is user | admin
	router.HandleFunc("/v1/key/{keyID}/{meta}", getKeyMetaHandler).Methods("GET")        // meta is ownwer | users | admins
	http.Handle("/", router)

	// run the web server
	http.ListenAndServe(":8080", nil)
}
