package db

import (
	"fmt"
	"librarymanagement/config"
	"log/slog"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect(DBConfig config.DB) *sqlx.DB {

	ConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DBConfig.Host,
		DBConfig.Port,
		DBConfig.User,
		DBConfig.Password,
		DBConfig.DbName,
		DBConfig.SSLMode)

	dbCon, err := sqlx.Connect("postgres", ConnStr)
	fmt.Println(ConnStr)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	dbCon.SetConnMaxIdleTime(
		time.Duration(10 * int(time.Minute)),
	)

	return dbCon

}

func ConnectDB() {
	conf := config.GetConfig()

	readDb = connect(conf.Db)
	slog.Info("Connected to read database")

	writeDb = connect(conf.Db)
	slog.Info("Connected to write database")
}

func CloseDB() {
	if err := readDb.Close(); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Disconnected from read database")

	if err := writeDb.Close(); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Disconnected from write database")
}
