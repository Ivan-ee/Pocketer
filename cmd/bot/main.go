package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegram-bot/pkg/config"
	"telegram-bot/pkg/repository"
	"telegram-bot/pkg/repository/boltdb"
	"telegram-bot/pkg/repository/server"
	"telegram-bot/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)

	if err != nil {
		log.Panic(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Panic(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	authServer := server.NewAuthServer(pocketClient, tokenRepository, cfg.BotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Panic(err)
		}
	}()

	if err := authServer.Start(); err != nil {
		log.Panic(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.BoltDBFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
