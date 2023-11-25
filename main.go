package main

import (
	"context"
	"flag"
	"log"

	clientsTelegram "github.com/empfaze/golang_bot/clients/telegram"
	"github.com/empfaze/golang_bot/consumer/event_consumer"
	eventsTelegram "github.com/empfaze/golang_bot/events/telegram"
	"github.com/empfaze/golang_bot/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

const (
	TG_BOT_HOST         = "api.telegram.org"
	SQLITE_STORAGE_PATH = "data/sqlite/storage.db"
	BATCH_SIZE          = 100
)

func main() {
	tgClient := clientsTelegram.New(TG_BOT_HOST, mustToken())

	storage, err := sqlite.New(SQLITE_STORAGE_PATH)
	if err != nil {
		log.Fatal("Couldn't connect to storage: ", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatal("Couldn't init storage: ", err)
	}

	eventsProcessor := eventsTelegram.New(tgClient, storage)

	log.Printf("Service has been started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, BATCH_SIZE)

	if err := consumer.Start(); err != nil {
		log.Fatal("Service has been stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-token-bot", "", "telegrem access token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}

	return *token
}
