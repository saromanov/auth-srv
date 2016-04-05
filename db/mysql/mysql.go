package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/auth-srv/db"
	account "github.com/micro/auth-srv/proto/account"
	oauth2 "github.com/micro/auth-srv/proto/oauth2"
)

var (
	Url = "haunted@root(127.0.0.1:3306)/auth"

	database string

	q = map[string]string{}

	// account queries
	accountQ = map[string]string{
		"delete":                "DELETE from %s.%s where id = ? limit 1",
		"create":                `INSERT into %s.%s (id, type, client_id, client_secret, salt, created, updated, metadata) values (?, ?, ?, ?, ?, ?, ?, ?)`,
		"update":                "UPDATE %s.%s set type = ?, client_secret = ?, salt = ?, updated = ?, metadata = ? where id = ?",
		"read":                  "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where id = ? limit 1",
		"search":                "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s limit ? offset ?",
		"searchClientId":        "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where client_id = ? limit ? offset ?",
		"searchType":            "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where type = ? limit ? offset ?",
		"searchClientIdAndType": "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s  where client_id = ? and type = ? limit ? offset ?",
	}

	// auth request queries
	authReqQ = map[string]string{
		"createReq":        "INSERT INTO %s.%s (id, response_type, client_id, scopes, state, redirect_uri, expires) values (?, ?, ?, ?, ?, ?, ?)",
		"deleteReq":        "DELETE FROM %s.%s where id = ? limit 1",
		"deleteReqExpired": "DELETE FROM %s.%s where expires <= ?",
		"readReq":          "SELECT id, response_type, client_id, scopes, state, redirect_uri, expires from %s.%s where id = ? limit 1",
	}

	tokenQ = map[string]string{
		"createToken": `INSERT INTO %s.%s (id, token_type, refresh_token, token_expires, refresh_expires, scopes, metadata, client_id, code) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"readToken": `SELECT id, token_type, refresh_token, token_expires, refresh_expires, scopes, metadata, client_id, code 
				from %s.%s where id = ? limit 1`,
		"updateToken":        "UPDATE %s.%s SET id = ?, refresh_token = ?, token_expires = ?, refresh_expires = ? where id = ? limit 1",
		"deleteToken":        "DELETE FROM %s.%s where id = ? limit 1",
		"deleteTokenExpired": "DELETE FROM %s.%s where refresh_expires <= ?",
		"readRefresh": `SELECT id, token_type, refresh_token, token_expires, refresh_expires, scopes, metadata, client_id, code 
				from %s.%s where refresh_token = ? limit 1`,
	}

	st = map[string]*sql.Stmt{}
)

type mysql struct {
	db *sql.DB
}

func init() {
	db.Register(new(mysql))
}

func (m *mysql) Init() error {
	var d *sql.DB
	var err error

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		return errors.New("Invalid database url")
	}

	if len(parts[1]) == 0 {
		return errors.New("Invalid database name")
	}

	url := parts[0]
	database := parts[1]

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		return err
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		return err
	}
	d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		return err
	}
	if _, err = d.Exec(accountSchema); err != nil {
		return err
	}
	if _, err = d.Exec(authReqSchema); err != nil {
		return err
	}
	if _, err = d.Exec(tokenSchema); err != nil {
		return err
	}

	for query, statement := range accountQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "accounts"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	for query, statement := range authReqQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "auth_requests"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	for query, statement := range tokenQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "tokens"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	m.db = d

	go m.run()

	return nil
}

func (m *mysql) run() {
	// mysql does not expire rows so we must do so.
	t := time.NewTicker(time.Minute * 10)

	for _ = range t.C {
		// expires field is actually set in the future
		// so we take anything older than now and destroy it
		expiry := time.Now().Unix()
		m.DeleteExpired(expiry)
	}
}

func (m *mysql) Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func (m *mysql) Create(acc *account.Record, salt string, secret string) error {
	b, err := json.Marshal(acc.Metadata)
	if err != nil {
		return err
	}
	meta := string(b)
	acc.Created = time.Now().Unix()
	acc.Updated = time.Now().Unix()
	_, err = st["create"].Exec(acc.Id, acc.Type, acc.ClientId, secret, salt, acc.Created, acc.Updated, meta)
	return err
}

func (m *mysql) Update(acc *account.Record, salt string, secret string) error {
	b, err := json.Marshal(acc.Metadata)
	if err != nil {
		return err
	}
	meta := string(b)
	acc.Updated = time.Now().Unix()
	_, err = st["update"].Exec(acc.Type, secret, salt, acc.Updated, meta, acc.Id)
	return err
}

func (m *mysql) Read(id string) (*account.Record, error) {
	acc := &account.Record{}

	r := st["read"].QueryRow(id)
	// we dont return salt or secret
	var salt, secret, meta string
	if err := r.Scan(&acc.Id, &acc.Type, &acc.ClientId, &secret, &salt, &acc.Created, &acc.Updated, &meta); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(meta), &acc.Metadata); err != nil {
		return nil, err
	}

	return acc, nil
}

func (m *mysql) Search(clientId, typ string, limit, offset int64) ([]*account.Record, error) {
	var r *sql.Rows
	var err error

	if len(clientId) > 0 && len(typ) > 0 {
		r, err = st["searchClientIdAndType"].Query(clientId, typ, limit, offset)
	} else if len(clientId) > 0 {
		r, err = st["searchClientId"].Query(clientId, limit, offset)
	} else if len(typ) > 0 {
		r, err = st["searchType"].Query(typ, limit, offset)
	} else {
		r, err = st["search"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var accounts []*account.Record

	for r.Next() {
		acc := &account.Record{}
		// we dont return salt or secret
		var salt, secret, meta string
		if err := r.Scan(&acc.Id, &acc.Type, &acc.ClientId, &secret, &salt, &acc.Created, &acc.Updated, &meta); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		if err := json.Unmarshal([]byte(meta), &acc.Metadata); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)

	}
	if r.Err() != nil {
		return nil, err
	}

	return accounts, nil
}

func (m *mysql) SaltAndSecret(id string) (string, string, error) {
	acc := &account.Record{}
	r := st["read"].QueryRow(id)
	var salt, secret, meta string
	if err := r.Scan(&acc.Id, &acc.Type, &acc.ClientId, &secret, &salt, &acc.Created, &acc.Updated, &meta); err != nil {
		if err == sql.ErrNoRows {
			return "", "", db.ErrNotFound
		}
		return "", "", err
	}
	return salt, secret, nil
}

// oauth2

func (m *mysql) CreateRequest(id string, req *oauth2.AuthorizeRequest) error {
	// expire in the future
	expiry := time.Now().Add(db.RequestExpiry).Unix()

	b, err := json.Marshal(req.Scopes)
	if err != nil {
		return err
	}

	_, err = st["createReq"].Exec(id, req.ResponseType, req.ClientId, string(b), req.State, req.RedirectUri, expiry)
	return err
}

func (m *mysql) DeleteRequest(id string) error {
	_, err := st["deleteReq"].Exec(id)
	return err
}

func (m *mysql) ReadRequest(id string) (*oauth2.AuthorizeRequest, error) {
	req := &oauth2.AuthorizeRequest{}

	r := st["readReq"].QueryRow(id)
	var scopes string
	var expires int64
	if err := r.Scan(&id, &req.ResponseType, &req.ClientId, &scopes, &req.State, &req.RedirectUri, &expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	// has it expired? if so, we're not giving this back
	if d := time.Now().Unix() - expires; d > 0 {
		return nil, db.ErrNotFound
	}

	if err := json.Unmarshal([]byte(scopes), &req.Scopes); err != nil {
		return nil, err
	}

	return req, nil
}

func (m *mysql) DeleteExpired(t int64) error {
	var gerr error
	_, err := st["deleteReqExpired"].Exec(t)
	if err != nil {
		gerr = err
	}
	_, err = st["deleteTokenExpired"].Exec(t)
	if err != nil {
		gerr = err
	}
	return gerr
}

func (m *mysql) ReadToken(id string) (*oauth2.Token, string, error) {
	tok := &oauth2.Token{}

	r := st["readToken"].QueryRow(id)
	var clientId, code, scopes, metadata string
	var refreshExpires int64

	if err := r.Scan(
		&tok.AccessToken,
		&tok.TokenType,
		&tok.RefreshToken,
		&tok.ExpiresAt,
		&refreshExpires,
		&scopes,
		&metadata,
		&clientId,
		&code,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", db.ErrNotFound
		}
		return nil, "", err
	}

	// has it expired? if so, we're not giving this back
	if d := time.Now().Unix() - refreshExpires; d > 0 {
		return nil, "", db.ErrNotFound
	}

	if err := json.Unmarshal([]byte(scopes), &tok.Scopes); err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal([]byte(metadata), &tok.Metadata); err != nil {
		return nil, "", err
	}

	return tok, clientId, nil
}

func (m *mysql) ReadRefresh(id string) (*oauth2.Token, string, error) {
	tok := &oauth2.Token{}

	r := st["readRefresh"].QueryRow(id)
	var clientId, code, scopes, metadata string
	var refreshExpires int64

	if err := r.Scan(
		&tok.AccessToken,
		&tok.TokenType,
		&tok.RefreshToken,
		&tok.ExpiresAt,
		&refreshExpires,
		&scopes,
		&metadata,
		&clientId,
		&code,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", db.ErrNotFound
		}
		return nil, "", err
	}

	// has it expired? if so, we're not giving this back
	if d := time.Now().Unix() - refreshExpires; d > 0 {
		return nil, "", db.ErrNotFound
	}

	if err := json.Unmarshal([]byte(scopes), &tok.Scopes); err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal([]byte(metadata), &tok.Metadata); err != nil {
		return nil, "", err
	}

	return tok, clientId, nil
}

func (m *mysql) CreateToken(t *oauth2.Token, clientId, code string) error {
	// expire in the future
	refreshExpiry := time.Now().Add(db.RefreshExpiry).Unix()

	bscope, err := json.Marshal(t.Scopes)
	if err != nil {
		return err
	}

	bmeta, err := json.Marshal(t.Metadata)
	if err != nil {
		return err
	}

	_, err = st["createToken"].Exec(
		t.AccessToken,
		t.TokenType,
		t.RefreshToken,
		t.ExpiresAt,
		refreshExpiry,
		string(bscope),
		string(bmeta),
		clientId,
		code,
	)

	return err
}

func (m *mysql) UpdateToken(id string, t *oauth2.Token) error {
	// expire in the future
	refreshExpiry := time.Now().Add(db.RefreshExpiry).Unix()

	_, err := st["updateToken"].Exec(
		t.AccessToken,
		t.RefreshToken,
		t.ExpiresAt,
		refreshExpiry,
		id,
	)

	return err
}

func (m *mysql) DeleteToken(id string) error {
	_, err := st["deleteToken"].Exec(id)
	return err
}
