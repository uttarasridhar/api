// Package tweets provides functionality around reading and writing tweets.
package tweets

type EmojisController interface {
	EmojiResults() ([]EmojiCount, error)
}

func NewEmojisController(db DB) EmojisController {
	return &emojisController{
		db: db,
	}
}

type emojisController struct {
	db DB
}

func (ec *emojisController) EmojiResults() ([]EmojiCount, error) {
	return ec.db.EmojiResults()
}
