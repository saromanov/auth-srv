# Auth Server

Auth server is an authentication and authorization microservice. It's used to authenticate both users and services. 
It also provides a mechanism for managing role based authorization.

Auth server currently implement Oauth2.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/auth-srv
	auth-srv --database_url="root:root@tcp(192.168.99.100:3306)/auth"
	```

	OR as a docker container

	```shell
	docker run microhq/auth-srv --database_url="root:root@tcp(192.168.99.100:3306)/auth" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Auth server implements the following RPC Methods

Account
- Read
- Create
- Update
- Delete
- Search

Oauth2
- Authorize
- Token
- Revoke

### Account.Create

```shell
micro query go.micro.srv.auth Account.Create '{"account": {"type": "user", "client_id": "asim", "client_secret": "foobar"}}'
```

### Account.Search

```shell
micro query go.micro.srv.auth Account.Search
{
	"accounts": [
		{
			"client_id": "asim",
			"created": 1.452816108e+09,
			"id": "2c02eea6-bb1b-11e5-9f39-68a86d0d36b6",
			"type": "user",
			"updated": 1.452816108e+09
		}
	]
}
```

### Oauth2.Authorize

Authorization Code Flow

```shell
micro query go.micro.srv.auth Oauth2.Authorize '{"response_type": "code", "client_id": "asim", "state": "mystatetoken", "redirect_uri": "https://foo.bar.com"}'
{
	"code": "cJMKdcx7iwAyhBLzNpmWQsSxpJOnuztB",
	"state": "mystatetoken"
}
```

### Oauth2.Token

Get Token

```shell
icro query go.micro.srv.auth Oauth2.Token '{"client_id": "asim", "client_secret": "foobar", "code": "cJMKdcx7iwAyhBLzNpmWQsSxpJOnuztB", "grant_type": "authorization_code", "redirect_uri": "https://foo.bar.com"}'
{
	"token": {
		"access_token": "V2swWmtsRm50WEtKSDhXSEtFdVlCNUo1WG5iTk9BYjh1dUVnT0JlOW9DS2FjWFg3c1FCaHBDbWFpaUhtQVUxUw==",
		"expires_at": 1.452819823e+09,
		"refresh_token": "OEZJUXBtdnNlTHNIWkhkRkQ4bTJFZkNNYlN6d0RQa2N6dkNwcDY1MkFCY0F5THdPZEFjdzB0a0JzNHpXYlJ4Ng==",
		"scopes": [
			"micro"
		],
		"token_type": "bearer"
	}
}
```

### Oauth2.Revoke

```shell
micro query go.micro.srv.auth Oauth2.Revoke '{"access_token": "V2swWmtsRm50WEtKSDhXSEtFdVlCNUo1WG5iTk9BYjh1dUVnT0JlOW9DS2FjWFg3c1FCaHBDbWFpaUhtQVUxUw=="}'
```
