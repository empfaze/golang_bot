package storage

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/empfaze/golang_bot/utils"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

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
