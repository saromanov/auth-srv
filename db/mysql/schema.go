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

	authReqSchema = `CREATE TABLE IF NOT EXISTS auth_requests (
id varchar(32) primary key,
response_type varchar(32),
client_id varchar(255),
scopes text,
state varchar(255),
redirect_uri text,
expires int(11),
index(expires));`

	tokenSchema = `CREATE TABLE IF NOT EXISTS tokens (
id varchar(255) primary key,
token_type varchar(16),
refresh_token varchar(255),
token_expires int(11),
refresh_expires int(11),
scopes text,
metadata text,
client_id varchar(255),
code varchar(32),
index(client_id),
index(token_expires),
index(refresh_expires),
unique(refresh_token));`
)
