package main

import (
	"flag"
	"log"

	"github.com/empfaze/golang_bot/clients/telegram"
)

const tgBotHost = "api.telegram.org"

func main() {
	token := mustToken()
	tgClient := telegram.New(tgBotHost, token)
}

func mustToken() string {
	token := flag.String("tg-token-bot", "", "telegrem access token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}

	return *token
}
