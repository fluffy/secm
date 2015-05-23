package main

import (
        "database/sql"
	_ "github.com/lib/pq"
        "log"
	"os"
)

var db *sql.DB

func setupDatabase( pgPassword string ) {
        var err error

        // set up DB
        db, err := sql.Open("postgres", "password='"+pgPassword+"' user=postgres dbname=postgres host=162.209.75.246 port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

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

func main() {
	var pgPassword string = os.Getenv("SECM_DB_SECRET")
	setupDatabase(pgPassword)
        defer db.Close()

}
