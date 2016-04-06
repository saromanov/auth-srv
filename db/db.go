package db

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	account "github.com/saromanov/auth-srv/proto/account"
	oauth2 "github.com/saromanov/auth-srv/proto/oauth2"
)

type DB interface {
	Init() error
	Account
	Oauth2
}

type Account interface {
	Read(id string) (*account.Record, error)
	Delete(id string) error
	Create(acc *account.Record, salt, secret string) error
	Update(acc *account.Record, salt, secret string) error
	Search(clientId, typ string, limit, offset int64) ([]*account.Record, error)
	SaltAndSecret(id string) (string, string, error)
}

type Oauth2 interface {
	ReadRequest(id string) (*oauth2.AuthorizeRequest, error)
	CreateRequest(id string, req *oauth2.AuthorizeRequest) error
	DeleteRequest(id string) error

	ReadToken(accessToken string) (*oauth2.Token, string, error)
	ReadRefresh(refreshToken string) (*oauth2.Token, string, error)
	CreateToken(token *oauth2.Token, clientId string, code string) error
	UpdateToken(accessToken string, token *oauth2.Token) error
	DeleteToken(accessToken string) error
}

var (
	db DB

	ErrNotFound = errors.New("not found")

	// we tick on this basis and kill anything older than this period
	RequestExpiry = time.Minute * 10
	RefreshExpiry = time.Hour * 24 * 14
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

func Token() string {
	return base64.StdEncoding.EncodeToString([]byte(random(64)))
}

// Request Code
func Code() string {
	return random(32)
}

// Some random salt
// bcrypt gives us one but we want ours too plus in binary
// fixed value so you need password, salt and binary to even
// begin hacking this.
func Salt() string {
	return random(16)
}

func Register(backend DB) {
	db = backend
}

func Init() error {
	return db.Init()
}

// account
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

// oauth2
func CreateRequest(id string, req *oauth2.AuthorizeRequest) error {
	return db.CreateRequest(id, req)
}

func DeleteRequest(id string) error {
	return db.DeleteRequest(id)
}

func ReadRequest(id string) (*oauth2.AuthorizeRequest, error) {
	return db.ReadRequest(id)
}

func ReadToken(accessToken string) (*oauth2.Token, string, error) {
	return db.ReadToken(accessToken)
}

func ReadRefresh(refreshToken string) (*oauth2.Token, string, error) {
	return db.ReadRefresh(refreshToken)
}

func CreateToken(t *oauth2.Token, clientId, code string) error {
	return db.CreateToken(t, clientId, code)

}

func UpdateToken(accessToken string, t *oauth2.Token) error {
	return db.UpdateToken(accessToken, t)

}

func DeleteToken(accessToken string) error {
	return db.DeleteToken(accessToken)
}
