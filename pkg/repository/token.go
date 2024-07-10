package repository

type Bucket string

const (
	AccessToken  Bucket = "access-tokens"
	RequestToken Bucket = "request-tokens"
)

type TokenRepository interface {
	Save(bucket Bucket, token string, chatId int64) error
	Get(bucket Bucket, chatId int64) (string, error)
}
