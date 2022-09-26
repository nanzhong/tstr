package server

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	identityv1 "github.com/nanzhong/tstr/api/identity/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestIdentityServer(t *testing.T) (*IdentityServer, *db.MockQuerier) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &IdentityServer{
		dbQuerier: mockQuerier,
	}, mockQuerier
}

func TestIdentityServer_Identity(t *testing.T) {
	server, mockQuerier := newTestIdentityServer(t)

	tokenString := "token"
	tokenHash := auth.HashToken(tokenString)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(auth.MDKeyAuth, "bearer "+tokenString))

	token := NewAccessTokenBuilder().Build()
	tokenID, err := uuid.Parse(token.Id)
	if err != nil {
		t.Skip("unable to parse the token id", token.Id)
	}

	mockQuerier.EXPECT().AuthAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenHash).Return(db.AuthAccessTokenRow{
		ID:                 tokenID,
		Name:               token.Name,
		NamespaceSelectors: token.NamespaceSelectors,
		Scopes:             []string{"admin"},
		IssuedAt:           sql.NullTime{Valid: true, Time: token.IssuedAt.AsTime()},
		ExpiresAt:          sql.NullTime{Valid: true, Time: token.ExpiresAt.AsTime()},
	}, nil)
	mockQuerier.EXPECT().ListAllNamespaces(ctx, gomock.AssignableToTypeOf(server.pgxPool)).Return([]string{"ns-0", "ns-1"}, nil)

	res, err := server.Identity(ctx, &identityv1.IdentityRequest{})
	require.NoError(t, err)
	assert.Equal(t, &identityv1.IdentityResponse{
		AccessToken:          token,
		AccessibleNamespaces: []string{"ns-0"},
	}, res)
}
