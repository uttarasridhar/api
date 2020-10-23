
// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"database/sql"
  "fmt"
)

const (
	// DBName is the name of the database for the tweets API.
	DBName = "tweets"
)

// DB is the interface for all the operations allowed on tweets.
type DB interface {
	StoreTweets(id, username, tweet_content, metadata string) error
  StoreEmojis(id, emoji string) error
	EmojiResults() ([]EmojiCount, error)
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
func (db *sqlDB) StoreTweets(id, username, tweet_content, metadata string) error {
	_, err := db.conn.Exec(`INSERT INTO tweets (id, username, tweet_content, metadata) VALUES ($1, $2, $3, $4)`, id, username, tweet_content, metadata)
	if err != nil {
		return fmt.Errorf("tweets: store tweet %s for tweet id %s and username %s: %v", tweet_content, id, username)
	}
	return nil
}

func (db *sqlDB) StoreEmojis(id, emoji string) error {
	_, err := db.conn.Exec(`INSERT INTO emojis (id, emoji) VALUES ($1, $2)`, id, emoji)
	if err != nil {
		return fmt.Errorf("emojis: store emoji for tweet id %s and emoji %s: %v", id, emoji)
	}
	return nil
}

// EmojiCount is a pair of a emoji and the count of occurrences of emoji.
type EmojiCount struct {
	Emoji string `json:"emoji"`
	Count  int    `json:"count"`
}

// EmojiResults returns the pair of emojis and counts.
func (db *sqlDB) EmojiResults() ([]EmojiCount, error) {
	rows, err := db.conn.Query(`SELECT emoji, COUNT(id) AS count FROM emojis GROUP BY emoji`)
	if err != nil {
		return nil, fmt.Errorf("emojis: retrieve emoji results: %v", err)
	}
	defer rows.Close()

	var results []EmojiCount
	for rows.Next() {
		var ec EmojiCount
		if err := rows.Scan(&ec.Emoji, &ec.Count); err != nil {
			return nil, fmt.Errorf("vote: scan row to emoji count pair: %v", err)
		}
		results = append(results, ec)
	}
	return results, nil
}

// CreateTweetsTableIfNotExist creates the "tweets" table if it does not exist already.
func CreateTweetsTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tweets (id VARCHAR(255) NOT NULL UNIQUE, username VARCHAR(255) NOT NULL, tweet_content VARCHAR(255) NOT NULL, metadata VARCHAR(255))`)
	if err != nil {
		return fmt.Errorf(`tweet: create "tweets" table: %v\n`, err)
	}
	return nil
}

// CreateEmojisTableIfNotExist creates the "tweets" table if it does not exist already.
func CreateEmojisTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS emojis (id VARCHAR(255), emoji VARCHAR(255), PRIMARY KEY (id, emoji))`)
	if err != nil {
		return fmt.Errorf(`tweet: create "emojis" table: %v\n`, err)
	}
	return nil
}
