package sqlcstorage

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/AlexeyInc/hw-otus/hw12_13_14_15_calendar/configs"
	_ "github.com/lib/pq"
)

var testQueries *Queries

var (
	configPath      = "./configs/config.toml"
	configLocalPath = "../../../../configs/config.toml"
)

func TestMain(m *testing.M) {
	config, err := configs.NewConfig(configLocalPath)
	if errors.Is(err, os.ErrNotExist) {
		config, err = configs.NewConfig(configPath)
	}
	if err != nil {
		log.Fatalln("can't read config file: " + err.Error())
	}
	testDB, err := sql.Open(config.Storage.Driver, config.Storage.Source)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	code := m.Run()

	deleteAllTestEvents()

	os.Exit(code)
}
