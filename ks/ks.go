package main

import (
        "database/sql"
	_ "github.com/lib/pq"
        "log"
	"fmt"
	"os"
)

var db *sql.DB

func setupDatabase( pgPassword string ) { // todo pass in hostname, port, username
        var err error

        // set up DB
        db, err = sql.Open("postgres", "password='"+pgPassword+"' user=postgres dbname=postgres host=162.209.75.246 port=5432 sslmode=disable")
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Setup Database\n")
}


func getKey( keyID, userID int64 ) {
	stmt, err := db.Prepare(
		"SELECT keys.kVal  FROM keyUsers JOIN keys ON  keys.kID = keyUsers.kID WHERE keyUsers.uID = $2 AND keyUsers.kID = $1")
	if err != nil {
		log.Println("sql fatal error in getKey prep")
		log.Fatal(err)
	}
	var keyVal string
	err = stmt.QueryRow(keyID,userID).Scan( &keyVal ) 
	if err != nil {
		log.Println("sql fatal error in getKey querry")
		log.Fatal(err)
	}
	fmt.Println(keyVal)
}


func main() {
	var err error
	
	var pgPassword string = os.Getenv("SECM_DB_SECRET")
	setupDatabase(pgPassword)
        defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	getKey( 101, 1 )
}
