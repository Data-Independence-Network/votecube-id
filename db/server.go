package db

import (
	"context"
	"database/sql"
	"votecube-id/models"

	"github.com/volatiletech/sqlboiler/boil"
)

var dBase *sql.DB

func SetupDb() *sql.DB {
	db, err := sql.Open("postgres", `postgresql://root@localhost:26257/votecube?sslmode=disable`)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	dBase = db

	return db
}

func SaveUser(u *models.User) error {
	return u.Insert(context.Background(), dBase, boil.Infer())
}
