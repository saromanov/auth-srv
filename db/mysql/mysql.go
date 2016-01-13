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
	"github.com/micro/auth-srv/proto/account"
)

var (
	Url = "root@tcp(127.0.0.1:3306)/auth"

	database string

	q = map[string]string{}

	accountQ = map[string]string{
		"delete":                "DELETE from %s.%s where id = ?",
		"create":                `INSERT into %s.%s (id, type, client_id, client_secret, salt, created, updated, metadata) values (?, ?, ?, ?, ?, ?, ?, ?)`,
		"update":                "UPDATE %s.%s set type = ?, client_secret = ?, salt = ?, updated = ?, metadata = ? where id = ?",
		"read":                  "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where id = ? limit 1",
		"search":                "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s limit ? offset ?",
		"searchClientId":        "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where client_id = ? limit ? offset ?",
		"searchType":            "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s where type = ? limit ? offset ?",
		"searchClientIdAndType": "SELECT id, type, client_id, client_secret, salt, created, updated, metadata from %s.%s  where client_id = ? and type = ? limit ? offset ?",
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

	for query, statement := range accountQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "accounts"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	m.db = d

	return nil
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
			return nil, errors.New("not found")
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
				return nil, errors.New("not found")
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
			return "", "", errors.New("not found")
		}
		return "", "", err
	}
	return salt, secret, nil
}
