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
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func newTestAdminServer(t *testing.T) (*AdminServer, *db.MockQuerier) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &AdminServer{
		dbQuerier: mockQuerier,
	}, mockQuerier
}

func TestAdminServer_IssueAccessToken(t *testing.T) {
	tests := []struct {
		name              string
		token             *commonv1.AccessToken
		responseCode      codes.Code
		errMsg            string
		mockQuerierReturn func(*commonv1.AccessToken) (db.IssueAccessTokenRow, error)
	}{
		{
			name:         "issue access token using valid token data request",
			token:        newAccessTokenBuilder().build(),
			responseCode: codes.OK,
			errMsg:       "",
			mockQuerierReturn: func(token *commonv1.AccessToken) (db.IssueAccessTokenRow, error) {
				tokenID, _ := uuid.Parse(token.Id)
				return db.IssueAccessTokenRow{
					ID:                 tokenID,
					Name:               token.Name,
					NamespaceSelectors: token.NamespaceSelectors,
					Scopes:             []string{"admin"},
					IssuedAt:           types.FromProtoTimestampAsNullTime(token.IssuedAt),
					ExpiresAt:          types.FromProtoTimestampAsNullTime(token.ExpiresAt),
				}, nil
			},
		},
		{
			name:         "db query fails when issuing access token",
			token:        newAccessTokenBuilder().build(),
			responseCode: codes.Internal,
			errMsg:       "failed to issue access token",
			mockQuerierReturn: func(token *commonv1.AccessToken) (db.IssueAccessTokenRow, error) {
				return db.IssueAccessTokenRow{}, errors.New("Dummy Error")
			},
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuerier.EXPECT().IssueAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), gomock.Any()).Return(tt.mockQuerierReturn(tt.token))

			request := adminv1.IssueAccessTokenRequest{
				Name:               tt.token.Name,
				NamespaceSelectors: tt.token.NamespaceSelectors,
				Scopes:             tt.token.Scopes,
				ValidDuration:      durationpb.New(tt.token.ExpiresAt.AsTime().Sub(tt.token.IssuedAt.AsTime())),
			}
			res, err := server.IssueAccessToken(ctx, &request)

			if res != nil {
				require.NoError(t, err)
				tt.token.Token = res.AccessToken.Token
				assert.Equal(t, &adminv1.IssueAccessTokenResponse{AccessToken: tt.token}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, er.Code(), tt.responseCode)
					assert.Equal(t, er.Message(), tt.errMsg)
				}
			}
		})
	}
}

func TestAdminServer_GetAccessToken(t *testing.T) {
	tests := []struct {
		name              string
		token             *commonv1.AccessToken
		responseCode      codes.Code
		errMsg            string
		mockQuerierReturn func(*commonv1.AccessToken) (db.GetAccessTokenRow, error)
	}{
		{
			name:         "get access token details using valid token uuid in the request",
			token:        newAccessTokenBuilder().withRevokedAt().build(),
			responseCode: codes.OK,
			errMsg:       "",
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
			name:              "get access token details using invalid token uuid in the request",
			token:             newAccessTokenBuilder().withRevokedAt().withId("invalid").build(),
			responseCode:      codes.InvalidArgument,
			errMsg:            "invalid access token id",
			mockQuerierReturn: nil,
		},
		{
			name:         "db query fails when getting access token",
			token:        newAccessTokenBuilder().build(),
			responseCode: codes.Internal,
			errMsg:       "failed to get access token",
			mockQuerierReturn: func(token *commonv1.AccessToken) (db.GetAccessTokenRow, error) {
				return db.GetAccessTokenRow{}, errors.New("Dummy Error")
			},
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	ctx := context.Background()

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
					assert.Equal(t, er.Code(), tt.responseCode)
					assert.Equal(t, er.Message(), tt.errMsg)
				}
			}
		})
	}
}

func TestAdminServer_ListAccessTokens(t *testing.T) {
	tests := []struct {
		name              string
		token             *commonv1.AccessToken
		responseCode      codes.Code
		errMsg            string
		mockQuerierReturn func(*commonv1.AccessToken) ([]db.ListAccessTokensRow, error)
	}{
		{
			name:         "valid access token available when listing access tokens",
			token:        newAccessTokenBuilder().withRevokedAt().build(),
			responseCode: codes.OK,
			errMsg:       "",
			mockQuerierReturn: func(token *commonv1.AccessToken) ([]db.ListAccessTokensRow, error) {
				tokenID, _ := uuid.Parse(token.Id)
				return []db.ListAccessTokensRow{{
					ID:                 tokenID,
					Name:               token.Name,
					NamespaceSelectors: token.NamespaceSelectors,
					Scopes:             []string{"admin"},
					IssuedAt:           types.FromProtoTimestampAsNullTime(token.IssuedAt),
					ExpiresAt:          types.FromProtoTimestampAsNullTime(token.ExpiresAt),
					RevokedAt:          types.FromProtoTimestampAsNullTime(token.RevokedAt),
				}}, nil
			},
		},
		{
			name:         "fail query for list access tokens",
			token:        newAccessTokenBuilder().build(),
			responseCode: codes.Internal,
			errMsg:       "failed to list access tokens",
			mockQuerierReturn: func(token *commonv1.AccessToken) ([]db.ListAccessTokensRow, error) {
				return []db.ListAccessTokensRow{}, errors.New("Dummy Error")
			},
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := db.ListAccessTokensParams{IncludeExpired: true, IncludeRevoked: true}
			mockQuerier.EXPECT().ListAccessTokens(ctx, gomock.AssignableToTypeOf(server.pgxPool), params).Return(tt.mockQuerierReturn(tt.token))

			request := &adminv1.ListAccessTokensRequest{IncludeExpired: params.IncludeExpired, IncludeRevoked: params.IncludeRevoked}
			res, err := server.ListAccessTokens(ctx, request)

			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &adminv1.ListAccessTokensResponse{AccessTokens: []*commonv1.AccessToken{tt.token}}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, er.Code(), tt.responseCode)
					assert.Equal(t, er.Message(), tt.errMsg)
				}
			}
		})
	}
}

func TestAdminServer_RevokeAccessToken(t *testing.T) {
	tests := []struct {
		name              string
		token             string
		responseCode      codes.Code
		errMsg            string
		mockQuerierReturn func() error
	}{
		{
			name:              "revoke access token using valid token uuid in the request",
			token:             uuid.New().String(),
			responseCode:      codes.OK,
			errMsg:            "",
			mockQuerierReturn: func() error { return nil },
		},
		{
			name:              "revoke access token using invalid token uuid in the request",
			token:             "invalid",
			responseCode:      codes.InvalidArgument,
			errMsg:            "invalid access token id",
			mockQuerierReturn: nil,
		},
		{
			name:              "fail query for revoke access tokens",
			token:             uuid.New().String(),
			responseCode:      codes.Internal,
			errMsg:            "failed to revoke access tokens",
			mockQuerierReturn: func() error { return errors.New("Dummy Error") },
		},
	}

	server, mockQuerier := newTestAdminServer(t)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockQuerierReturn != nil {
				tokenID, _ := uuid.Parse(tt.token)
				mockQuerier.EXPECT().RevokeAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenID).Return(tt.mockQuerierReturn())
			}

			request := &adminv1.RevokeAccessTokenRequest{Id: tt.token}
			res, err := server.RevokeAccessToken(ctx, request)

			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &adminv1.RevokeAccessTokenResponse{}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, er.Code(), tt.responseCode)
					assert.Equal(t, er.Message(), tt.errMsg)
				}
			}
		})
	}
}
