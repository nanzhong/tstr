package server

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newTestAdminServer(t *testing.T) (*AdminServer, *db.MockQuerier) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &AdminServer{
		dbQuerier: mockQuerier,
	}, mockQuerier
}

func TestAdminServer_GetAccessToken(t *testing.T) {
	tests := []struct {
		name              string
		token             *commonv1.AccessToken
		errCode           codes.Code
		errMsg            string
		mockQuerierReturn func(*commonv1.AccessToken) (db.GetAccessTokenRow, error)
	}{
		{
			name:    "valid access token uuid in the request",
			token:   newAccessTokenBuilder().withRevokedAt().build(),
			errCode: codes.OK,
			errMsg:  "",
			mockQuerierReturn: func(token *commonv1.AccessToken) (db.GetAccessTokenRow, error) {
				tokenID, _ := uuid.Parse(token.Id)
				return db.GetAccessTokenRow{
					ID:                 tokenID,
					Name:               token.Name,
					NamespaceSelectors: token.NamespaceSelectors,
					Scopes:             []string{"admin"},
					IssuedAt:           types.FromProtoTimestampAsNullTime(token.IssuedAt),
					ExpiresAt:          types.FromProtoTimestampAsNullTime(token.ExpiresAt),
					RevokedAt:          types.FromProtoTimestampAsNullTime(token.RevokedAt),
				}, nil
			},
		},
		{
			name:              "invalid access token uuid in the request",
			token:             newAccessTokenBuilder().withRevokedAt().withId("invalid").build(),
			errCode:           codes.InvalidArgument,
			errMsg:            "invalid access token id",
			mockQuerierReturn: nil,
		},
		{
			name:    "fail query for access token",
			token:   newAccessTokenBuilder().build(),
			errCode: codes.Internal,
			errMsg:  "failed to get access token",
			mockQuerierReturn: func(token *commonv1.AccessToken) (db.GetAccessTokenRow, error) {
				return db.GetAccessTokenRow{}, errors.New("Dummy Error")
			},
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	tokenString := "token"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(auth.MDKeyAuth, "bearer "+tokenString))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockQuerierReturn != nil {
				tokenID, _ := uuid.Parse(tt.token.Id)
				mockQuerier.EXPECT().GetAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenID).Return(tt.mockQuerierReturn(tt.token))
			}

			request := &adminv1.GetAccessTokenRequest{Id: tt.token.Id}
			res, err := server.GetAccessToken(ctx, request)

			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &adminv1.GetAccessTokenResponse{AccessToken: tt.token}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, er.Code(), tt.errCode)
					assert.Equal(t, er.Message(), tt.errMsg)
				}
			}
		})
	}
}
