// Package api starts the api server.
package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/uttarasridhar/api/postgres"
	"github.com/uttarasridhar/api/server"
	"github.com/uttarasridhar/api/tweets"

	"github.com/gorilla/mux"
)

// Run starts the server.
func Run() error {
	addr := flag.String("addr", ":8080", "port to listen on")
	flag.Parse()

	secret := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := json.Unmarshal([]byte(os.Getenv("RDS_SECRET")), &secret); err != nil {
		return fmt.Errorf("api: unmarshal rds secret: %v", err)
	}

	conn, close, err := postgres.Connect(
		os.Getenv("RDS_ENDPOINT"),
		postgres.Port,
		secret.Username,
		secret.Password,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)
	if err != nil {
		return fmt.Errorf("api: connect to postgres db: %v", err)
	}
	defer close()
	if err := tweets.CreateTweetsTableIfNotExist(conn); err != nil {
		return fmt.Errorf("api: create table: %v", err)
	}
	if err := tweets.CreateEmojisTableIfNotExist(conn); err != nil {
		return fmt.Errorf("api: create table: %v", err)
	}

	db := tweets.NewSQLDB(conn)

	tweets := tweets.NewTweetsController(db)
	emojis := tweets.NewEmojisController(db)

	s := http.Server{
		Addr: *addr,
		Handler: &server.Server{
			Router: mux.NewRouter(),
			Tweets: tweets,
			Emojis: emojis,
		},
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("listen on port %s\n", *addr)
	return s.ListenAndServe()
}
