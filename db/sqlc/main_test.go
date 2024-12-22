package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

var testQueries *db.Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
