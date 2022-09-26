package server

import (
	"time"

	"github.com/google/uuid"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/grpc/types"
)

type AccessTokenBuilder struct {
	AccessToken *commonv1.AccessToken
}

func NewAccessTokenBuilder() *AccessTokenBuilder {
	return &AccessTokenBuilder{
		&commonv1.AccessToken{
			Id:                 uuid.New().String(),
			Name:               "name",
			NamespaceSelectors: []string{"ns-0"},
			Scopes:             []commonv1.AccessToken_Scope{commonv1.AccessToken_ADMIN},
			IssuedAt:           types.ToProtoTimestamp(time.Now()),
			ExpiresAt:          types.ToProtoTimestamp(time.Now().Add(time.Hour)),
		},
	}
}

func (b *AccessTokenBuilder) WithId(id string) *AccessTokenBuilder {
	b.AccessToken.Id = id
	return b
}

func (b *AccessTokenBuilder) WithName(name string) *AccessTokenBuilder {
	b.AccessToken.Name = name
	return b
}

func (b *AccessTokenBuilder) WithNamespaceSelectors(namespaceSelectors []string) *AccessTokenBuilder {
	b.AccessToken.NamespaceSelectors = namespaceSelectors
	return b
}

func (b *AccessTokenBuilder) WithScopes(scopes []commonv1.AccessToken_Scope) *AccessTokenBuilder {
	b.AccessToken.Scopes = scopes
	return b
}

func (b *AccessTokenBuilder) WithRevokedAt() *AccessTokenBuilder {
	b.AccessToken.RevokedAt = types.ToProtoTimestamp(time.Now().Add(time.Hour))
	return b
}

func (b *AccessTokenBuilder) Build() *commonv1.AccessToken {
	return b.AccessToken
}
