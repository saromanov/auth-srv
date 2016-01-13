package handler

import (
	"github.com/micro/auth-srv/proto/oauth2"

	"golang.org/x/net/context"
)

type Oauth2 struct{}

func (o *Oauth2) Authorize(ctx context.Context, req *oauth2.AuthorizeRequest, rsp *oauth2.AuthorizeResponse) error {
	return nil
}

func (o *Oauth2) Token(ctx context.Context, req *oauth2.TokenRequest, rsp *oauth2.TokenResponse) error {
	return nil
}
