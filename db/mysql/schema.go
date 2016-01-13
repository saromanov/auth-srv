package mysql

var (
	accountSchema = `CREATE TABLE IF NOT EXISTS accounts (
id varchar(36) primary key,
type varchar(64),
client_id varchar(255),
client_secret text,
salt varchar(16),
created integer,
updated integer,
metadata text,
unique (client_id));`
)
