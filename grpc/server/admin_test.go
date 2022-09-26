package server

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
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
		name    string
		token   *commonv1.AccessToken
		errCode codes.Code
		errMsg  string
	}{
		{
			"valid acces token uuid in the request",
			NewAccessTokenBuilder().WithRevokedAt().Build(),
			codes.OK,
			"",
		},
		{
			"invalid acces token uuid in the request, multiple scopes",
			NewAccessTokenBuilder().WithRevokedAt().WithId("invalid").Build(),
			codes.InvalidArgument,
			"invalid access token id",
		},
		{
			"fail query for access token",
			nil,
			codes.Internal,
			"failed to get access token",
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	tokenString := "token"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(auth.MDKeyAuth, "bearer "+tokenString))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request *adminv1.GetAccessTokenRequest
			if tt.token != nil {
				tokenID, err := uuid.Parse(tt.token.Id)
				if err == nil {
					mockQuerier.EXPECT().GetAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenID).Return(db.GetAccessTokenRow{
						ID:                 tokenID,
						Name:               tt.token.Name,
						NamespaceSelectors: tt.token.NamespaceSelectors,
						Scopes:             []string{"admin"},
						IssuedAt:           sql.NullTime{Valid: true, Time: tt.token.IssuedAt.AsTime()},
						ExpiresAt:          sql.NullTime{Valid: true, Time: tt.token.ExpiresAt.AsTime()},
						RevokedAt:          sql.NullTime{Valid: true, Time: tt.token.RevokedAt.AsTime()},
					}, nil)
				}
				request = &adminv1.GetAccessTokenRequest{Id: tt.token.Id}
			} else {
				tokenID := uuid.New()
				mockQuerier.EXPECT().GetAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenID).Return(db.GetAccessTokenRow{}, errors.New("Dummy Error"))
				request = &adminv1.GetAccessTokenRequest{Id: tokenID.String()}
			}

			res, err := server.GetAccessToken(ctx, request)
			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &adminv1.GetAccessTokenResponse{AccessToken: tt.token}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", tt.errCode, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}