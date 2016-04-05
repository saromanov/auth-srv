package restful

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"github.com/saromanov/auth-srv/proto/account"
	"github.com/saromanov/auth-srv/proto/oauth2"

	"golang.org/x/net/context"
)

type User struct{}

var (
	cl account.AccountClient
	auth oauth2.Oauth2Client
)

func Create(w http.ResponseWriter, r *http.Request) {
	var req account.CreateRequest
	var err error
	//Trying to decode request body
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//var resp account.UpdateNameResponse
	resp, err := cl.Create(context.TODO(), &req)

	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(result))
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var req oauth2.AuthorizeRequest
	var err error
	//Trying to decode request body
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//var resp account.UpdateNameResponse
	resp, err := auth.Authorize(context.TODO(), &req)

	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(result))
}

func Token(w http.ResponseWriter, r *http.Request) {
	var req oauth2.TokenRequest
	var err error
	//Trying to decode request body
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//var resp account.UpdateNameResponse
	resp, err := auth.Token(context.TODO(), &req)

	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(resp)
	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(result))
}

func InitRestful() {
	service := web.NewService(
		web.Name("go.micro.srv.auth"),
	)

	service.Init()

	// setup Greeter Server Client
	cl = account.NewAccountClient("go.micro.srv.auth", client.DefaultClient)

	http.HandleFunc("/auth/authorize", Auth)
	http.HandleFunc("/auth/token", Token)
	http.ListenAndServe(":8082", nil)
}