// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"log"
)

type TweetsController interface {
	StoreTweets(id, username, tweet_content, metadata string) error
}

func NewTweetsController(db DB) TweetsController {
	return &tweetsController{
		db: db,
	}
}

type tweetsController struct {
	db DB
}

func (tc *tweetsController) StoreTweets(id, username, tweet_content, metadata string) error {
	err := tc.db.StoreTweets(id, username, tweet_content, metadata)
	if err != nil {
		return err
	}
	emojis, err := tc.convertTweetToEmojisList(tweet_content)
	if err != nil {
		// TODO: We need some kind of recon that reconciles missing data between tables.
		return err
	}
	for _, emoji := range emojis {
		err := tc.db.StoreEmojis(id, emoji)
		if err != nil {
			log.Printf("ERROR: server: store id=%s emoji=%s: %v\n", id, emoji, err)
		}
	}
	return nil
}

func (tc *tweetsController) convertTweetToEmojisList(tweet_content string) ([]string, error) {
	// TODO: PH, thank you
	return []string {}, nil
}
