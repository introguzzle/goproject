package persistence

import (
	"database/sql"
	"fmt"
	"goproject/src/env"
)

var (
	database *sql.DB
)

var (
	host       string
	port       string
	user       string
	name       string
	password   string
	connection string
)

func buildConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
}

func Close() {
	if database != nil {
		err := database.Close()
		if err != nil {
			return
		}
	}
}

func Connect() *sql.DB {
	db, err := sql.Open(connection, buildConnectionString())

	if err != nil {
		panic(err)
	}

	return db
}

func init() {
	host = env.Get("DB_HOST").Value
	port = env.Get("DB_PORT").Value
	user = env.Get("DB_USER").Value
	name = env.Get("DB_NAME").Value
	password = env.Get("DB_PASSWORD").Value
	connection = env.Get("DB_CONNECTION").Value
}

func Ping() error {
	return Connect().Ping()
}
