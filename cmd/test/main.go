package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rpoletaev/huskyjam/internal"
)

func main(){
	db, err := sqlx.Connect("postgres", "host=localhost user=postgres dbname=huskyjam password=suppersecret sslmode=disable")
    if err != nil {
        log.Fatalln(err)
	}

	log.Println("connected")
	acc := &internal.Account{
		Email: "new@acc.zone",
		Pass: "supper",
	}
	if _, err := db.NamedExec("INSERT INTO accounts (email, pass) VALUES (:email, :pass)", acc); err != nil {
		log.Fatalln(err)
	}

	log.Println("created")
}