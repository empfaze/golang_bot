package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/empfaze/golang_bot/lib/storage"
	"github.com/empfaze/golang_bot/utils"
)

const (
	RND_CMD   = "/rnd"
	HELP_CMD  = "/help"
	START_CMD = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("Got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(ctx, text, chatID, username)
	}

	switch text {
	case RND_CMD:
		return p.sendRandom(ctx, chatID, username)
	case HELP_CMD:
		return p.SendHelp(chatID)
	case START_CMD:
		return p.SendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, pageUrl string, chatID int, username string) error {
	page := &storage.Page{
		URL:      pageUrl,
		Username: username,
	}

	isExist, err := p.storage.IsExist(ctx, page)
	if err != nil {
		return utils.WrapError("Couldn't save page: ", err)
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return utils.WrapError("Couldn't save page: ", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return utils.WrapError("Couldn't save page: ", err)
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) error {
	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return utils.WrapError("Couldn't send random: ", err)
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	result, err := url.Parse(text)

	return err == nil && result.Host != ""
}
