package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/empfaze/golang_bot/lib/storage"
	"github.com/empfaze/golang_bot/utils"
)

const (
	defaultPermission = 0774
	errMessage        = "An error occured while saving: "
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	path := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(path, defaultPermission); err != nil {
		return utils.WrapError(errMessage, err)
	}

	name, err := fileName(page)
	if err != nil {
		return utils.WrapError("Error while hashing data: ", err)
	}

	filePath := filepath.Join(path, name)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(len(files))

	file := files[number]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, utils.WrapError("Can't decode page", err)
	}

	return &page, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return utils.WrapError("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return utils.WrapError(msg, err)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, utils.WrapError("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, utils.WrapError(msg, err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
