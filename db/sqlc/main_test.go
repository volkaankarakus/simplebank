package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var error error
	testDB, error = sql.Open(dbDriver, dbSource)
	if error != nil {
		log.Fatal("cannot connect to database:", error)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
