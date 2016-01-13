package db

import (
	"crypto/rand"
	"errors"

	"github.com/micro/auth-srv/proto/account"
)

type DB interface {
	Init() error
	Read(id string) (*account.Record, error)
	Delete(id string) error
	Create(acc *account.Record, salt, secret string) error
	Update(acc *account.Record, salt, secret string) error
	Search(clientId, typ string, limit, offset int64) ([]*account.Record, error)
	SaltAndSecret(id string) (string, string, error)
}

var (
	db DB

	ErrNotFound = errors.New("not found")
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
	return "ughwhy?!!!"
}

func Salt() string {
	return random(16)
}

func Register(backend DB) {
	db = backend
}

func Init() error {
	return db.Init()
}

func Read(id string) (*account.Record, error) {
	return db.Read(id)
}

func Create(ch *account.Record, salt, secret string) error {
	return db.Create(ch, salt, secret)
}

func Update(ch *account.Record, salt, secret string) error {
	return db.Update(ch, salt, secret)
}

func Delete(id string) error {
	return db.Delete(id)
}

func Search(clientId, typ string, limit, offset int64) ([]*account.Record, error) {
	return db.Search(clientId, typ, limit, offset)
}

func SaltAndSecret(id string) (string, string, error) {
	return db.SaltAndSecret(id)
}
