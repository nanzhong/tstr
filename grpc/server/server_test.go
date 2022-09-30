package server

import (
	"time"

	"github.com/google/uuid"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/grpc/types"
)

type accessTokenBuilder struct {
	accessToken *commonv1.AccessToken
}

func NewAccessTokenBuilder() *accessTokenBuilder {
	return &accessTokenBuilder{
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

func (b *accessTokenBuilder) WithId(id string) *accessTokenBuilder {
	b.accessToken.Id = id
	return b
}

func (b *accessTokenBuilder) WithName(name string) *accessTokenBuilder {
	b.accessToken.Name = name
	return b
}

func (b *accessTokenBuilder) WithNamespaceSelectors(namespaceSelectors []string) *accessTokenBuilder {
	b.accessToken.NamespaceSelectors = namespaceSelectors
	return b
}

func (b *accessTokenBuilder) WithScopes(scopes []commonv1.AccessToken_Scope) *accessTokenBuilder {
	b.accessToken.Scopes = scopes
	return b
}

func (b *accessTokenBuilder) WithRevokedAt() *accessTokenBuilder {
	b.accessToken.RevokedAt = types.ToProtoTimestamp(time.Now().Add(time.Hour))
	return b
}

func (b *accessTokenBuilder) Build() *commonv1.AccessToken {
	return b.accessToken
}
