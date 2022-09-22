package server

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	identityv1 "github.com/nanzhong/tstr/api/identity/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/types"
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

	tokenID := uuid.New()
	token := &commonv1.AccessToken{
		Id:                 tokenID.String(),
		Name:               "name",
		NamespaceSelectors: []string{"ns-0"},
		Scopes:             []commonv1.AccessToken_Scope{commonv1.AccessToken_ADMIN},
		IssuedAt:           types.ToProtoTimestamp(time.Now()),
		ExpiresAt:          types.ToProtoTimestamp(time.Now().Add(time.Hour)),
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
