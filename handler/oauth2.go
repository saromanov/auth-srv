package handler

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/micro/auth-srv/db"
	"github.com/micro/auth-srv/proto/oauth2"
	"github.com/micro/go-micro/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var (
	DefaultScope = "micro"
)

type Oauth2 struct{}

func authClient(clientId, clientSecret string) error {
	acc, err := db.Search(clientId, "", 1, 0)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth", "server_error")
	}

	if len(acc) == 0 {
		return errors.BadRequest("go.micro.srv.auth", "invalid_request")
	}

	// check the secret
	salt, secret, err := db.SaltAndSecret(acc[0].Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth", "server_error")
	}

	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.auth", "server_error")
	}

	// does it match?
	if err := bcrypt.CompareHashAndPassword(s, []byte(x+salt+clientSecret)); err != nil {
		return errors.BadRequest("go.micro.srv.auth", "access_denied")
	}

	return nil
}

func (o *Oauth2) Authorize(ctx context.Context, req *oauth2.AuthorizeRequest, rsp *oauth2.AuthorizeResponse) error {
	// We may actually need to authenticate who can make this request.
	// How should we do that?

	switch req.ResponseType {
	// requesting authorization code
	case "code":
		// check client id exists
		if len(req.ClientId) == 0 {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// if redirect uri exists and is not tls lets bail
		if len(req.RedirectUri) > 0 && !strings.HasPrefix(req.RedirectUri, "https://") {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// use default scope
		if len(req.Scopes) == 0 {
			req.Scopes = append(req.Scopes, DefaultScope)
		}

		// generate code
		code := db.Code()

		// store request; expire in 10 mins
		if err := db.CreateRequest(code, req); err != nil {
			return errors.InternalServerError("go.micro.srv.auth", "server_error")
		}

		// respond

		rsp.Code = code
		rsp.State = req.State

		// we're done?!

	// implicit token request
	case "token":
		// to be implemented
		return errors.BadRequest("go.micro.srv.auth", "unsupported_response_type")
	default:
		return errors.BadRequest("go.micro.srv.auth", "unsupported_response_type")
	}

	return nil
}

func (o *Oauth2) Token(ctx context.Context, req *oauth2.TokenRequest, rsp *oauth2.TokenResponse) error {
	// We may actually need to authenticate who can make this request.
	// How should we do that?

	// TODO: track multiple attempts for the same authorization code and revoke all other tokens

	// supported grant types
	switch req.GrantType {
	case "authorization_code":
		// validate inputs
		if len(req.Code) == 0 || len(req.ClientId) == 0 || len(req.ClientSecret) == 0 {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// read authorization request
		storedReq, err := db.ReadRequest(req.Code)
		if err != nil {
			if err == db.ErrNotFound {
				return errors.BadRequest("go.micro.srv.auth", "invalid_request")
			} else {
				return errors.InternalServerError("go.micro.srv.auth", "server_error")
			}
		}
		// we have a request, is it the same client_id?
		if req.ClientId != storedReq.ClientId {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// is it the same redirect uri?
		if req.RedirectUri != storedReq.RedirectUri {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// auth
		if err := authClient(req.ClientId, req.ClientSecret); err != nil {
			return err
		}

		// ok successful auth

		// generate a token; tokens are basically opaque strings
		// we just hand back something random and base64 encoded
		token := &oauth2.Token{
			AccessToken:  db.Token(),
			TokenType:    "bearer",
			RefreshToken: db.Token(),
			ExpiresAt:    time.Now().Add(time.Hour).Unix(),
			Scopes:       storedReq.Scopes,
			Metadata:     req.Metadata,
		}

		// store the token against the client
		if err := db.CreateToken(token, req.ClientId, req.Code); err != nil {
			return errors.InternalServerError("go.micro.srv.auth", "server_error")
		}

		// return token
		rsp.Token = token
	case "client_credentials":
		// validate inputs
		if len(req.ClientId) == 0 || len(req.ClientSecret) == 0 {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// auth
		if err := authClient(req.ClientId, req.ClientSecret); err != nil {
			return err
		}

		// ok successful auth

		// generate a token; tokens are basically opaque strings
		// we just hand back something random and base64 encoded
		token := &oauth2.Token{
			AccessToken:  db.Token(),
			TokenType:    "bearer",
			RefreshToken: db.Token(),
			ExpiresAt:    time.Now().Add(time.Hour).Unix(),
			Scopes:       req.Scopes,
			Metadata:     req.Metadata,
		}

		// store the token against the client
		if err := db.CreateToken(token, req.ClientId, ""); err != nil {
			return errors.InternalServerError("go.micro.srv.auth", "server_error")
		}

		// return token
		rsp.Token = token
	case "refresh_token":
		// validate inputs
		if len(req.RefreshToken) == 0 || len(req.ClientId) == 0 || len(req.ClientSecret) == 0 {
			return errors.BadRequest("go.micro.srv.auth", "invalid_request")
		}

		// auth client
		if err := authClient(req.ClientId, req.ClientSecret); err != nil {
			return err
		}

		// get existing token
		token, clientId, err := db.ReadRefresh(req.RefreshToken)
		if err != nil {
			if err == db.ErrNotFound {
				return errors.BadRequest("go.micro.srv.auth", "invalid_request")
			} else {
				return errors.InternalServerError("go.micro.srv.auth", "server_error")
			}
		}

		// client id does not match for refresh token
		if clientId != req.ClientId {
			return errors.BadRequest("go.micro.srv.auth", "access_denied")
		}

		// so we have a token, we now need to refresh

		id := token.AccessToken
		token.AccessToken = db.Token()
		token.RefreshToken = db.Token()
		token.ExpiresAt = time.Now().Add(time.Hour).Unix()

		// Update with new access token, refresh token and expiry
		// refresh expiry is set by db
		if err := db.UpdateToken(id, token); err != nil {
			return errors.InternalServerError("go.micro.srv.auth", "server_error")
		}

		rsp.Token = token
	default:
		return errors.BadRequest("go.micro.srv.auth", "unsupported_grant_type")
	}

	// Delete the authorization request when we successfully create a token
	if err := db.DeleteRequest(req.Code); err != nil {
		return errors.InternalServerError("go.micro.srv.auth", "server_error")
	}

	return nil
}

func (o *Oauth2) Revoke(ctx context.Context, req *oauth2.RevokeRequest, rsp *oauth2.RevokeResponse) error {
	// Who should be allowed to do this?

	if len(req.RefreshToken) > 0 {
		token, _, err := db.ReadRefresh(req.RefreshToken)
		if err != nil {
			if err == db.ErrNotFound {
				return errors.BadRequest("go.micro.srv.auth", "invalid_request")
			} else {
				return errors.InternalServerError("go.micro.srv.auth", "server_error")
			}
		}

		if err := db.DeleteToken(req.AccessToken); err != nil {
			return errors.InternalServerError("go.micro.srv.auth", "server_error")
		}

		req.AccessToken = token.AccessToken
	}

	if len(req.AccessToken) == 0 {
		return errors.BadRequest("go.micro.srv.auth", "invalid_request")
	}

	if err := db.DeleteToken(req.AccessToken); err != nil {
		return errors.InternalServerError("go.micro.srv.auth", "server_error")
	}

	return nil
}
