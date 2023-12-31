package dbutil_test

import (
	"database/sql"
	"github.com/fortytw2/dockertest"
	_ "github.com/jinzhu/gorm/dialects/postgres" // DB adapter
	dbutil "github.com/nekizz/init-project/pkg/util/db"
	"testing"
)

func TestDatabase(t *testing.T) {
	container, err := dockertest.RunContainer("postgres:9.6", "5432", func(addr string) error {
		db, err := sql.Open("postgres", "postgres://postgres:postgres@"+addr+"?sslmode=disable")
		if err != nil {
			return err
		}

		return db.Ping()
	})
	defer container.Shutdown()
	if err != nil {
		t.Fatalf("could not start postgres, %s", err)
	}

	_, err = dbutil.New("postgres", "PSN", false)
	if err == nil {
		t.Error("Expected error")
	}

	_, err = dbutil.New("postgres", "postgres://postgres:postgres@localhost:1234/postgres?sslmode=disable", false)
	if err == nil {
		t.Error("Expected error")
	}

	dbLogTest, err := dbutil.New("postgres", "postgres://postgres:postgres@"+container.Addr+"/postgres?sslmode=disable", true)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}
	dbLogTest.Close()
}
