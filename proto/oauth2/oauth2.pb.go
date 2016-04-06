// Code generated by protoc-gen-go.
// source: oauth2.proto
// DO NOT EDIT!

/*
Package go_micro_srv_auth_oauth2 is a generated protocol buffer package.

It is generated from these files:
	oauth2.proto

It has these top-level messages:
	Token
	AuthorizeRequest
	AuthorizeResponse
	TokenRequest
	TokenResponse
	RevokeRequest
	RevokeResponse
	ValidateRequest
	ValidateResponse
*/
package go_micro_srv_auth_oauth2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Token struct {
	AccessToken  string   `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	TokenType    string   `protobuf:"bytes,2,opt,name=token_type,json=tokenType" json:"token_type,omitempty"`
	RefreshToken string   `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
	ExpiresAt    int64    `protobuf:"varint,4,opt,name=expires_at,json=expiresAt" json:"expires_at,omitempty"`
	Scopes       []string `protobuf:"bytes,5,rep,name=scopes" json:"scopes,omitempty"`
	// metadata associated with the token
	Metadata map[string]string `protobuf:"bytes,6,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Token) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type AuthorizeRequest struct {
	// code, token (not supported)
	ResponseType string   `protobuf:"bytes,1,opt,name=response_type,json=responseType" json:"response_type,omitempty"`
	ClientId     string   `protobuf:"bytes,2,opt,name=client_id,json=clientId" json:"client_id,omitempty"`
	Scopes       []string `protobuf:"bytes,3,rep,name=scopes" json:"scopes,omitempty"`
	State        string   `protobuf:"bytes,4,opt,name=state" json:"state,omitempty"`
	RedirectUri  string   `protobuf:"bytes,5,opt,name=redirect_uri,json=redirectUri" json:"redirect_uri,omitempty"`
}

func (m *AuthorizeRequest) Reset()                    { *m = AuthorizeRequest{} }
func (m *AuthorizeRequest) String() string            { return proto.CompactTextString(m) }
func (*AuthorizeRequest) ProtoMessage()               {}
func (*AuthorizeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type AuthorizeResponse struct {
	Code  string `protobuf:"bytes,1,opt,name=code" json:"code,omitempty"`
	State string `protobuf:"bytes,2,opt,name=state" json:"state,omitempty"`
	// implicit response
	Token *Token `protobuf:"bytes,3,opt,name=token" json:"token,omitempty"`
}

func (m *AuthorizeResponse) Reset()                    { *m = AuthorizeResponse{} }
func (m *AuthorizeResponse) String() string            { return proto.CompactTextString(m) }
func (*AuthorizeResponse) ProtoMessage()               {}
func (*AuthorizeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AuthorizeResponse) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

type TokenRequest struct {
	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret" json:"client_secret,omitempty"`
	Code         string `protobuf:"bytes,3,opt,name=code" json:"code,omitempty"`
	// password (not supported), client_credentials, authorization_code, refresh_token
	GrantType    string `protobuf:"bytes,4,opt,name=grant_type,json=grantType" json:"grant_type,omitempty"`
	RedirectUri  string `protobuf:"bytes,5,opt,name=redirect_uri,json=redirectUri" json:"redirect_uri,omitempty"`
	RefreshToken string `protobuf:"bytes,6,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
	// scopes can be added for client_credentials request
	Scopes []string `protobuf:"bytes,7,rep,name=scopes" json:"scopes,omitempty"`
	// metadata to be stored with a token that's generated
	Metadata map[string]string `protobuf:"bytes,8,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *TokenRequest) Reset()                    { *m = TokenRequest{} }
func (m *TokenRequest) String() string            { return proto.CompactTextString(m) }
func (*TokenRequest) ProtoMessage()               {}
func (*TokenRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *TokenRequest) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type TokenResponse struct {
	Token *Token `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *TokenResponse) Reset()                    { *m = TokenResponse{} }
func (m *TokenResponse) String() string            { return proto.CompactTextString(m) }
func (*TokenResponse) ProtoMessage()               {}
func (*TokenResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *TokenResponse) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

type RevokeRequest struct {
	// revoke access token
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
	// revoke via refresh token
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken" json:"refresh_token,omitempty"`
}

func (m *RevokeRequest) Reset()                    { *m = RevokeRequest{} }
func (m *RevokeRequest) String() string            { return proto.CompactTextString(m) }
func (*RevokeRequest) ProtoMessage()               {}
func (*RevokeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type RevokeResponse struct {
}

func (m *RevokeResponse) Reset()                    { *m = RevokeResponse{} }
func (m *RevokeResponse) String() string            { return proto.CompactTextString(m) }
func (*RevokeResponse) ProtoMessage()               {}
func (*RevokeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type ValidateRequest struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken" json:"access_token,omitempty"`
}

func (m *ValidateRequest) Reset()                    { *m = ValidateRequest{} }
func (m *ValidateRequest) String() string            { return proto.CompactTextString(m) }
func (*ValidateRequest) ProtoMessage()               {}
func (*ValidateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type ValidateResponse struct {
	Token  *Token `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
	Active bool   `protobuf:"varint,2,opt,name=active" json:"active,omitempty"`
}

func (m *ValidateResponse) Reset()                    { *m = ValidateResponse{} }
func (m *ValidateResponse) String() string            { return proto.CompactTextString(m) }
func (*ValidateResponse) ProtoMessage()               {}
func (*ValidateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ValidateResponse) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

func init() {
	proto.RegisterType((*Token)(nil), "go.micro.srv.auth.oauth2.Token")
	proto.RegisterType((*AuthorizeRequest)(nil), "go.micro.srv.auth.oauth2.AuthorizeRequest")
	proto.RegisterType((*AuthorizeResponse)(nil), "go.micro.srv.auth.oauth2.AuthorizeResponse")
	proto.RegisterType((*TokenRequest)(nil), "go.micro.srv.auth.oauth2.TokenRequest")
	proto.RegisterType((*TokenResponse)(nil), "go.micro.srv.auth.oauth2.TokenResponse")
	proto.RegisterType((*RevokeRequest)(nil), "go.micro.srv.auth.oauth2.RevokeRequest")
	proto.RegisterType((*RevokeResponse)(nil), "go.micro.srv.auth.oauth2.RevokeResponse")
	proto.RegisterType((*ValidateRequest)(nil), "go.micro.srv.auth.oauth2.ValidateRequest")
	proto.RegisterType((*ValidateResponse)(nil), "go.micro.srv.auth.oauth2.ValidateResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Oauth2 service

type Oauth2Client interface {
	Authorize(ctx context.Context, in *AuthorizeRequest, opts ...client.CallOption) (*AuthorizeResponse, error)
	Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error)
	Revoke(ctx context.Context, in *RevokeRequest, opts ...client.CallOption) (*RevokeResponse, error)
	Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error)
}

type oauth2Client struct {
	c           client.Client
	serviceName string
}

func NewOauth2Client(serviceName string, c client.Client) Oauth2Client {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.auth.oauth2"
	}
	return &oauth2Client{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *oauth2Client) Authorize(ctx context.Context, in *AuthorizeRequest, opts ...client.CallOption) (*AuthorizeResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Oauth2.Authorize", in)
	out := new(AuthorizeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauth2Client) Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Oauth2.Token", in)
	out := new(TokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauth2Client) Revoke(ctx context.Context, in *RevokeRequest, opts ...client.CallOption) (*RevokeResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Oauth2.Revoke", in)
	out := new(RevokeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauth2Client) Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Oauth2.Validate", in)
	out := new(ValidateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Oauth2 service

type Oauth2Handler interface {
	Authorize(context.Context, *AuthorizeRequest, *AuthorizeResponse) error
	Token(context.Context, *TokenRequest, *TokenResponse) error
	Revoke(context.Context, *RevokeRequest, *RevokeResponse) error
	Validate(context.Context, *ValidateRequest, *ValidateResponse) error
}

func RegisterOauth2Handler(s server.Server, hdlr Oauth2Handler) {
	s.Handle(s.NewHandler(&Oauth2{hdlr}))
}

type Oauth2 struct {
	Oauth2Handler
}

func (h *Oauth2) Authorize(ctx context.Context, in *AuthorizeRequest, out *AuthorizeResponse) error {
	return h.Oauth2Handler.Authorize(ctx, in, out)
}

func (h *Oauth2) Token(ctx context.Context, in *TokenRequest, out *TokenResponse) error {
	return h.Oauth2Handler.Token(ctx, in, out)
}

func (h *Oauth2) Revoke(ctx context.Context, in *RevokeRequest, out *RevokeResponse) error {
	return h.Oauth2Handler.Revoke(ctx, in, out)
}

func (h *Oauth2) Validate(ctx context.Context, in *ValidateRequest, out *ValidateResponse) error {
	return h.Oauth2Handler.Validate(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 592 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x18, 0x25, 0xc9, 0x1a, 0x9a, 0x6f, 0x2d, 0x14, 0x0b, 0xa1, 0xa8, 0x08, 0x31, 0x82, 0x04, 0x63,
	0x88, 0x5c, 0x94, 0x21, 0x21, 0xb8, 0xda, 0x05, 0x48, 0xbb, 0x40, 0xa0, 0x30, 0x40, 0x42, 0x42,
	0x95, 0x71, 0xbd, 0x2d, 0x5a, 0x57, 0x07, 0xdb, 0xad, 0x28, 0x4f, 0xc1, 0x4b, 0x70, 0xc9, 0xcb,
	0xf1, 0x04, 0x38, 0xb6, 0x93, 0x26, 0x2d, 0x25, 0x9d, 0x10, 0x77, 0xf1, 0xf1, 0xf7, 0x77, 0xce,
	0xf9, 0xea, 0x42, 0x87, 0xe1, 0xa9, 0x3c, 0x1d, 0xc4, 0x19, 0x67, 0x92, 0xa1, 0xf0, 0x84, 0xc5,
	0xe7, 0x29, 0xe1, 0x2c, 0x16, 0x7c, 0x16, 0xe7, 0x37, 0xb1, 0xb9, 0x8f, 0x7e, 0xba, 0xd0, 0x3a,
	0x62, 0x67, 0x74, 0x82, 0xee, 0x40, 0x07, 0x13, 0x42, 0x85, 0x18, 0xca, 0xfc, 0x1c, 0x3a, 0x3b,
	0xce, 0x6e, 0x90, 0x6c, 0x1b, 0xcc, 0x84, 0xdc, 0x02, 0xd0, 0x77, 0x43, 0x39, 0xcf, 0x68, 0xe8,
	0xea, 0x80, 0x40, 0x23, 0x47, 0x0a, 0x40, 0x77, 0xa1, 0xcb, 0xe9, 0x31, 0xa7, 0xe2, 0xd4, 0x96,
	0xf0, 0x74, 0x44, 0xc7, 0x82, 0x65, 0x0d, 0xfa, 0x35, 0x4b, 0x15, 0x30, 0xc4, 0x32, 0xdc, 0x52,
	0x11, 0x5e, 0x12, 0x58, 0xe4, 0x40, 0xa2, 0x1b, 0xe0, 0x0b, 0xc2, 0x32, 0x2a, 0xc2, 0xd6, 0x8e,
	0xa7, 0x92, 0xed, 0x09, 0x1d, 0x42, 0xfb, 0x9c, 0x4a, 0x3c, 0xc2, 0x12, 0x87, 0xbe, 0xba, 0xd9,
	0x1e, 0x3c, 0x8a, 0xd7, 0x91, 0x8a, 0x75, 0xa7, 0xf8, 0x95, 0x8d, 0x7f, 0x31, 0x91, 0x7c, 0x9e,
	0x94, 0xe9, 0xfd, 0xe7, 0xd0, 0xad, 0x5d, 0xa1, 0x1e, 0x78, 0x67, 0x74, 0x6e, 0x09, 0xe7, 0x9f,
	0xe8, 0x3a, 0xb4, 0x66, 0x78, 0x3c, 0x2d, 0x38, 0x9a, 0xc3, 0x33, 0xf7, 0xa9, 0x13, 0xfd, 0x70,
	0xa0, 0x77, 0xa0, 0x9a, 0x30, 0x9e, 0x7e, 0xa3, 0x09, 0xfd, 0x32, 0xa5, 0x42, 0x1a, 0xe2, 0x22,
	0x63, 0x13, 0x41, 0x8d, 0x34, 0x4e, 0x41, 0xdc, 0x80, 0x5a, 0x9d, 0x9b, 0x10, 0x90, 0x71, 0x4a,
	0x27, 0x72, 0x98, 0x8e, 0x6c, 0xdd, 0xb6, 0x01, 0x0e, 0x47, 0x15, 0xda, 0x5e, 0x8d, 0xb6, 0x1a,
	0x44, 0x48, 0x2c, 0xa9, 0x16, 0x4a, 0x0d, 0xa2, 0x0f, 0xb9, 0x55, 0x9c, 0x8e, 0x94, 0x62, 0x44,
	0x0e, 0xa7, 0x3c, 0x55, 0x52, 0x69, 0xab, 0x0a, 0xec, 0x1d, 0x4f, 0x23, 0x09, 0xd7, 0x2a, 0x63,
	0x9a, 0x31, 0x10, 0x82, 0x2d, 0xc2, 0x46, 0xc5, 0x78, 0xfa, 0x7b, 0xd1, 0xc1, 0xad, 0x76, 0x78,
	0x02, 0xad, 0x85, 0x85, 0xdb, 0x83, 0xdb, 0x0d, 0x5a, 0x27, 0x26, 0x3a, 0xfa, 0xe5, 0x42, 0xc7,
	0x00, 0x56, 0x99, 0x1a, 0x69, 0x67, 0x89, 0xb4, 0x92, 0xcd, 0x5e, 0x0a, 0x4a, 0x38, 0x95, 0x76,
	0x84, 0x8e, 0x01, 0xdf, 0x6a, 0xac, 0x9c, 0xd9, 0xab, 0xcc, 0xac, 0x76, 0xe8, 0x84, 0x63, 0x95,
	0xa7, 0xc5, 0x36, 0xd2, 0x04, 0x1a, 0xd1, 0x4a, 0x37, 0xcb, 0xb3, 0xba, 0xaa, 0xfe, 0x1f, 0x56,
	0x75, 0x61, 0xca, 0xe5, 0x9a, 0x29, 0x6f, 0x2a, 0xbb, 0xd8, 0xd6, 0xbb, 0xb8, 0xdf, 0xa4, 0x8f,
	0x91, 0xe3, 0xff, 0xac, 0xe4, 0x4b, 0xe8, 0xda, 0x26, 0xd6, 0xe6, 0xd2, 0x3c, 0xe7, 0x42, 0xe6,
	0x7d, 0x80, 0x6e, 0x42, 0x67, 0xea, 0xb3, 0x30, 0x6f, 0x83, 0x17, 0x61, 0x45, 0x47, 0x77, 0x55,
	0xc7, 0xa8, 0x07, 0x57, 0x8a, 0xc2, 0x66, 0xc2, 0x68, 0x1f, 0xae, 0xbe, 0xc7, 0xe3, 0x54, 0xf1,
	0xbd, 0x40, 0xb3, 0x08, 0x43, 0x6f, 0x91, 0xf5, 0x4f, 0x5c, 0x73, 0x6b, 0x31, 0x91, 0xe9, 0xcc,
	0xc8, 0xd9, 0x4e, 0xec, 0x69, 0xf0, 0xdd, 0x03, 0xff, 0xb5, 0x8e, 0x47, 0xc7, 0x10, 0x94, 0xbf,
	0x20, 0xb4, 0xb7, 0xbe, 0xee, 0xf2, 0x6b, 0xd0, 0x7f, 0xb8, 0x51, 0xac, 0x55, 0xe2, 0x12, 0xfa,
	0x58, 0x3c, 0xc0, 0xf7, 0x36, 0x5b, 0xa2, 0xfe, 0xfd, 0xc6, 0xb8, 0xb2, 0xf6, 0x27, 0xf0, 0x8d,
	0xf2, 0xe8, 0x2f, 0x49, 0x35, 0xd3, 0xfb, 0xbb, 0xcd, 0x81, 0x65, 0x79, 0x02, 0xed, 0xc2, 0x10,
	0xf4, 0x60, 0x7d, 0xde, 0x92, 0xd5, 0xfd, 0xbd, 0x4d, 0x42, 0x8b, 0x26, 0x9f, 0x7d, 0xfd, 0x17,
	0xf6, 0xf8, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x51, 0xd9, 0xbb, 0xc9, 0xd2, 0x06, 0x00, 0x00,
}
