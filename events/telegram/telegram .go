package telegram

import (
	"errors"

	"github.com/empfaze/golang_bot/clients/telegram"
	"github.com/empfaze/golang_bot/events"
	"github.com/empfaze/golang_bot/lib/storage"
	"github.com/empfaze/golang_bot/utils"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEvent = errors.New("Unknown event type")
	ErrUnknownMeta  = errors.New("Unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, utils.WrapError("Couldn't get events: ", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		result = append(result, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.MESSAGE:
		return p.processMessage(event)
	default:
		return utils.WrapError("Couldn't process the message: ", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return utils.WrapError("Couldn't process message: ", err)
	}

	if err != p.doCmd(event.Text, meta.ChatID, meta.Username) {
		return utils.WrapError("Couldn't process message: ", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, utils.WrapError("Couldn't get meta: ", ErrUnknownMeta)
	}

	return result, nil
}

func event(update telegram.Update) events.Event {
	updType := fetchType(update)
	updText := fetchText(update)

	result := events.Event{
		Type: updType,
		Text: updText,
	}

	if updType == events.MESSAGE {
		result.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return result
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.UNKNOWN
	}

	return events.MESSAGE
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
