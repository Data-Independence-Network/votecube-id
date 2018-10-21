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
	Email            string
	Password         string
	VerificationCode string
	Found            bool
	ResponseCh       chan bool
}

type UserState int

const (
	NewUser   UserState = 0
	OauthUser UserState = 1
	PwdUser   UserState = 2
)

var UserRequestCh = make(chan *LoginInfo)

var incomingLoginInfos = make([]*LoginInfo, 0, 1)

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
			email := loginInfo.Email
			emails = append(emails, email)
			loginInfoMap[email] = loginInfo
		}

		users, err := models.Users(Where("email in ?", emails)).All(context.Background(), dBase)

		if err != nil {
			for _, loginInfo := range loginInfosToProcess {
				loginInfo.ResponseCh <- false
			}
			log.Printf("DB error: %v\n", err)
			continue
		}

		for _, user := range users {
			loginInfo := loginInfoMap[user.Email]
			loginInfo.Found = true
			loginInfo.ResponseCh <- true
			close(loginInfo.ResponseCh)
		}

		for _, loginInfo := range loginInfosToProcess {
			if !loginInfo.Found {
				loginInfo.ResponseCh <- false
				close(loginInfo.ResponseCh)
			}
		}

	}
}

func SaveUser(u *models.User) error {
	return u.Insert(context.Background(), dBase, boil.Infer())
}
