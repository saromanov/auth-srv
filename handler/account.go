package handler

import (
	"encoding/base64"
	"strings"

	"github.com/saromanov/auth-srv/db"
	"github.com/micro/go-micro/errors"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	account "github.com/saromanov/auth-srv/proto/account"
)

const (
	x = "_micro_auth_"
)

type Account struct{}

func validateAccount(acc *account.Record, method string) error {
	if acc == nil {
		return errors.BadRequest("go.micro.srv.auth."+method, "invalid account")
	}

	if len(acc.Id) > 0 && uuid.Parse(acc.Id) == nil {
		return errors.BadRequest("go.micro.srv.auth."+method, "invalid id")
	}

	if len(acc.Type) == 0 {
		return errors.BadRequest("go.micro.srv.auth."+method, "type cannot be blank")
	}

	if len(acc.ClientId) == 0 {
		return errors.BadRequest("go.micro.srv.auth."+method, "client id cannot be blank")
	}

	if len(acc.ClientSecret) == 0 {
		return errors.BadRequest("go.micro.srv.auth."+method, "client secret cannot be blank")
	}

	return nil
}

func (s *Account) Create(ctx context.Context, req *account.CreateRequest, rsp *account.CreateResponse) error {
	// validate incoming
	if err := validateAccount(req.Account, "Create"); err != nil {
		return err
	}

	// set a uuid if we dont have one
	if len(req.Account.Id) == 0 {
		req.Account.Id = uuid.NewUUID().String()
	}

	// hash the pass
	salt := db.Salt()
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Account.ClientSecret), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	// to lower
	req.Account.ClientId = strings.ToLower(req.Account.ClientId)
	req.Account.Type = strings.ToLower(req.Account.Type)

	if err := db.Create(req.Account, salt, pp); err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Create", err.Error())
	}
	return nil
}

func (s *Account) Read(ctx context.Context, req *account.ReadRequest, rsp *account.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.auth.Read", "id cannot be blank")
	}

	acc, err := db.Read(req.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Read", err.Error())
	}

	rsp.Account = acc
	return nil
}

func (s *Account) Update(ctx context.Context, req *account.UpdateRequest, rsp *account.UpdateResponse) error {
	// validate incoming
	if err := validateAccount(req.Account, "Update"); err != nil {
		return err
	}

	// need an account id for update
	if len(req.Account.Id) == 0 {
		return errors.BadRequest("go.micro.srv.auth.Update", "invalid id")
	}

	// lookup the record and verify it's the same
	acc, err := db.Read(req.Account.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Update", err.Error())
	}

	// not the same client id
	if req.Account.ClientId != acc.ClientId {
		return errors.BadRequest("go.micro.srv.auth.Update", "invalid client id")
	}

	// hash the pass
	salt := db.Salt()
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Account.ClientSecret), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Update", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	// to lower
	req.Account.ClientId = strings.ToLower(req.Account.ClientId)
	req.Account.Type = strings.ToLower(req.Account.Type)

	// update
	if err := db.Update(req.Account, salt, pp); err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Update", err.Error())
	}

	return nil
}

func (s *Account) Delete(ctx context.Context, req *account.DeleteRequest, rsp *account.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.auth.Delete", "invalid id")
	}

	if err := db.Delete(req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.auth.Delete", err.Error())
	}
	return nil
}

func (s *Account) Search(ctx context.Context, req *account.SearchRequest, rsp *account.SearchResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	accounts, err := db.Search(req.ClientId, req.Type, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Accounts = accounts
	return nil
}
