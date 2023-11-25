package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/empfaze/golang_bot/utils"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, username string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

var ErrNoSavedPages error = errors.New("No saved pages")

func (p Page) Hash() (string, error) {
	hash := sha1.New()

	if _, err := io.WriteString(hash, p.URL); err != nil {
		return "", utils.WrapError("Couldn't write url: ", err)
	}

	if _, err := io.WriteString(hash, p.Username); err != nil {
		return "", utils.WrapError("Couldn't write username: ", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
