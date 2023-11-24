package main

import (
	"flag"
	"log"

	clientsTelegram "github.com/empfaze/golang_bot/clients/telegram"
	"github.com/empfaze/golang_bot/consumer/event_consumer"
	eventsTelegram "github.com/empfaze/golang_bot/events/telegram"
	"github.com/empfaze/golang_bot/lib/files"
)

const (
	TG_BOT_HOST  = "api.telegram.org"
	STORAGE_PATH = "storage"
	BATCH_SIZE   = 100
)

func main() {
	tgClient := clientsTelegram.New(TG_BOT_HOST, mustToken())
	eventsProcessor := eventsTelegram.New(tgClient, files.New(STORAGE_PATH))

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
