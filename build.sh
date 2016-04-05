#!/usr/bin/env bash

echo "Get dependencies"
go get ./...
echo "Build service"
go build -o ./bin/service main.go
echo "Run service"
./bin/service

Well, at least this is works locally. And already tried on the DO.

Start consumer service
$GOPATH/bin/consumer-srv

Start users service
$GOPATH/bin/user-srv

Start test service
$GOPATH/bin/test-srv

Create new users. Event sending to rabbit and message should be shown on consumer-srv
micro query go.micro.srv.user Account.Create '{"user":{"username": "foobar", "email": "foobar@foobar.com"}, "password": "123"}'

Login user.  Event sending to rabbit and message should be shown on consumer-srv
micro query go.micro.srv.user Account.Login '{"email": "foobar@foobar.com", "password": "123"}'

Change username via rest  Event sending to rabbit and message should be shown on consumer-srv
curl -H "Content-Type: application/json" -X POST -d '{"email": "foobar@foobar.com", "oldusername":"foobar", "username":"boom"}' http://localhost:8082/test/update