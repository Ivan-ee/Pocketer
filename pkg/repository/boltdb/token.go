package bolt

import (
	"errors"
	"github.com/boltdb/bolt"
	"strconv"
	"telegram-bot/pkg/repository"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(bucket repository.Bucket, token string, chatId int64) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatId), []byte(token))
	})
}

func (r *TokenRepository) Get(bucket repository.Bucket, chatId int64) (string, error) {
	var token string

	err := r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatId))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}

func intToBytes(value int64) []byte {
	return []byte(strconv.FormatInt(value, 10))
}
