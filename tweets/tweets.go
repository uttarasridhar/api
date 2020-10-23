// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	// DBName is the name of the database for the tweets API.
	DBName = "tweets"
)

// DB is the interface for all the operations allowed on tweets.
type DB interface {
	Store(id, username, tweet_content, metadata string) error
	EmojiResults() (results []EmojiCount, err error)
}

// EmojiCount is a pair of a emoji and the count of occurrences of emoji.
type EmojiCount struct {
	Emoji string `json:"emoji"`
	Count  int    `json:"count"`
}

// NewSQLDB creates a sql database to read and store tweets.
func NewSQLDB(db *sql.DB) DB {
	return &sqlDB{
		conn: db,
	}
}

type execQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type sqlDB struct {
	conn execQuerier
}

// Store a tweet in the database.
func (db *sqlDB) Store(id, username, tweet_content, metadata string) error {
	_, err := db.conn.Exec(`INSERT INTO tweets (id, username, tweet_content, metadata) VALUES ($1, $2, $3, $4)`, id, username, tweet_content, metadata)
	if err == nil {
		return nil
	}
	log.Printf("INFO: tweets: update tweet for tweet id %s and username %s\n", id, username)
	_, err = db.conn.Exec(`UPDATE tweets SET username = $1, tweet_content = $2, metadata = $3 WHERE id = $4`, username, tweet_content, metadata, id)
	if err != nil {
		return fmt.Errorf("tweets: store tweet %s for tweet id %s and username %s: %v", tweet_content, id, username)
	}
	return nil
}

// EmojiResults returns the pair of emojis and counts.
func (db *sqlDB) EmojiResults() ([]EmojiCount, error) {
	return nil, nil
}

// CreateTableIfNotExist creates the "tweets" table if it does not exist already.
func CreateTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tweets (id VARCHAR(255) NOT NULL UNIQUE, username VARCHAR(255) NOT NULL, tweet_content VARCHAR(255) NOT NULL, metadata VARCHAR(255))`)
	if err != nil {
		return fmt.Errorf(`tweet: create "tweets" table: %v\n`, err)
	}
	return nil
}
