package db

import (
	"context"
	"database/sql"
	"log"
	"time"
	"votecube-id/models"

	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

var dBase *sql.DB

type LoginInfo struct {
	email            string
	password         string
	verificationCode string
	found            bool
	responseCh       chan bool
}

var UserRequestCh = make(chan LoginInfo)

var incomingLoginInfos = make([]LoginInfo, 0, 1)

const SLEEP_TIME = 1000 * time.Millisecond

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

	go collectLoginInfos()

	go processLoginInfos()

	return db
}

func collectLoginInfos() {
	for {
		incomingLoginInfos = append(incomingLoginInfos, <-UserRequestCh)
	}
}

func processLoginInfos() {
	for {
		time.Sleep(SLEEP_TIME)

		var loginInfosToProcess = incomingLoginInfos
		incomingLoginInfos = make([]LoginInfo, 0, 1)
		if len(loginInfosToProcess) == 0 {
			log.Printf("No record to process.\n")
			continue
		}
		log.Printf("Processing records: %v\n", len(loginInfosToProcess))

		loginInfoMap := make(map[string]LoginInfo)
		emails := make([]string, len(loginInfosToProcess))

		for _, loginInfo := range loginInfosToProcess {
			email := loginInfo.email
			emails = append(emails, email)
			loginInfoMap[email] = loginInfo
		}

		users, err := models.Users(Where("email in ?", emails)).All(context.Background(), dBase)

		if err != nil {
			for _, loginInfo := range loginInfosToProcess {
				loginInfo.responseCh <- false
			}
			log.Printf("DB error: %v\n", err)
			continue
		}

		for _, user := range users {
			loginInfo := loginInfoMap[user.Email]
			loginInfo.found = true
			loginInfo.responseCh <- true
		}

		for _, loginInfo := range loginInfosToProcess {
			if !loginInfo.found {
				loginInfo.responseCh <- false
			}
		}

	}
}

func SaveUser(u *models.User) error {
	return u.Insert(context.Background(), dBase, boil.Infer())
}
