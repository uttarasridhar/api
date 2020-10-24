// Package tweets provides functionality around reading and writing tweets.
package tweets

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	rangeMin1 = 127744
	rangeMax1 = 129750
	rangeMin2 = 126980
	rangeMax2 = 127569
	rangeMin3 = 169
	rangeMax3 = 174
	rangeMin4 = 8205
	rangeMax4 = 12953
)

type TweetsController interface {
	StoreTweets(id, username, tweetContent, metadata string) error
}

func NewTweetsController(db DB) TweetsController {
	return &tweetsController{
		db: db,
	}
}

type tweetsController struct {
	db DB
}

func (tc *tweetsController) StoreTweets(id, username, tweetContent, metadata string) error {
	err := tc.db.StoreTweets(id, username, tweetContent, metadata)
	if err != nil {
		return err
	}
	emojis, err := tc.convertTweetToEmojisList(tweetContent)
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

// convertTweetToEmojisList extracts emoji from a tweet and store in decimal.
func (tc *tweetsController) convertTweetToEmojisList(tweetContent string) ([]string, error) {
	var emoji []string
	r, err := regexp.Compile(`\\U\w{8}`)
	if err != nil {
		return nil, err
	}
	for _, e := range r.FindAllString(tweetContent, -1) {
		dec, err := hexUnicodeToDecUnicode(e)
		if err != nil {
			return nil, err
		}
		emoji = append(emoji, dec)
	}
	replaced := string(r.ReplaceAll([]byte(tweetContent), []byte("")))
	for _, r := range replaced {
		if isEmoji(r) {
			emoji = append(emoji, fmt.Sprint(r))
		}
	}
	return emoji, nil
}

// Example: \U0001f9c3 -> 129475
func hexUnicodeToDecUnicode(s string) (string, error) {
	hex := strings.TrimPrefix(s, `\U`)
	dec, err := strconv.ParseInt(hex, 16, 32)
	if err == nil {
		return fmt.Sprint(dec), nil
	}
	return "", err
}

func isEmoji(r rune) bool {
	code := int(r)
	switch {
	case code >= rangeMin1 && code <= rangeMax1:
		return true
	case code >= rangeMin2 && code <= rangeMax2:
		return true
	case code >= rangeMin3 && code <= rangeMax3:
		return true
	case code >= rangeMin4 && code <= rangeMax4:
		return true
	default:
		return false
	}
}
