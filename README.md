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


