package server

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	adminv1 "github.com/nanzhong/tstr/api/admin/v1"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminServer struct {
	adminv1.UnimplementedAdminServiceServer

	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func NewAdminServer(pgxPool *pgxpool.Pool) adminv1.AdminServiceServer {
	return &AdminServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
	}
}

func (s *AdminServer) IssueAccessToken(ctx context.Context, req *adminv1.IssueAccessTokenRequest) (*adminv1.IssueAccessTokenResponse, error) {
	var expiresAt sql.NullTime
	if req.ValidDuration != nil {
		expiresAt.Valid = true
		expiresAt.Time = time.Now().UTC().Add(req.ValidDuration.AsDuration())
	}

	tokenData := make([]byte, 32)
	_, err := rand.Read(tokenData)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token value")
		return nil, status.Error(codes.Internal, "failed to issue access token")
	}
	token := hex.EncodeToString(tokenData)

	tokenHashBytes := sha512.Sum512([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashBytes[:])

	textScopes := make([]string, len(req.Scopes))
	for i, s := range types.FromAccessTokenScopes(req.Scopes) {
		textScopes[i] = string(s)
	}
	issuedToken, err := s.dbQuerier.IssueAccessToken(ctx, s.pgxPool, db.IssueAccessTokenParams{
		Name:      req.Name,
		TokenHash: tokenHash,
		Scopes:    textScopes,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create access token in db")
		return nil, status.Error(codes.Internal, "failed to issue access token")
	}

	return &adminv1.IssueAccessTokenResponse{
		AccessToken: &commonv1.AccessToken{
			Id:   issuedToken.ID.String(),
			Name: issuedToken.Name,
			// TODO return actual scopes inserted, pending sqlc bug fix re enum arrays
			// Scopes:    types.ToAccessTokenScopes(issuedToken.Scopes),
			Scopes:    req.Scopes,
			Token:     token,
			IssuedAt:  types.ToProtoTimestamp(issuedToken.IssuedAt),
			ExpiresAt: types.ToProtoTimestamp(issuedToken.ExpiresAt),
		},
	}, nil
}

func (s *AdminServer) GetAccessToken(ctx context.Context, req *adminv1.GetAccessTokenRequest) (*adminv1.GetAccessTokenResponse, error) {
	tokenID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid access token id")
	}

	token, err := s.dbQuerier.GetAccessToken(ctx, s.pgxPool, tokenID)
	if err != nil {
		log.Error().Err(err).Msg("failed to query for access token")
		return nil, status.Error(codes.Internal, "failed to get access token")
	}

	return &adminv1.GetAccessTokenResponse{
		AccessToken: &commonv1.AccessToken{
			Id:        token.ID.String(),
			Name:      token.Name,
			Scopes:    types.ToAccessTokenScopes(token.Scopes),
			IssuedAt:  types.ToProtoTimestamp(token.IssuedAt),
			ExpiresAt: types.ToProtoTimestamp(token.ExpiresAt),
			RevokedAt: types.ToProtoTimestamp(token.RevokedAt),
		},
	}, nil
}

func (s *AdminServer) ListAccessTokens(ctx context.Context, req *adminv1.ListAccessTokensRequest) (*adminv1.ListAccessTokensResponse, error) {
	tokens, err := s.dbQuerier.ListAccessTokens(ctx, s.pgxPool, db.ListAccessTokensParams{
		IncludeExpired: req.IncludeExpired,
		IncludeRevoked: req.IncludeRevoked,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to query for access tokens")
		return nil, status.Error(codes.Internal, "failed to list access tokens")
	}

	res := &adminv1.ListAccessTokensResponse{}
	for _, token := range tokens {
		res.AccessTokens = append(res.AccessTokens, &commonv1.AccessToken{
			Id:        token.ID.String(),
			Name:      token.Name,
			Scopes:    types.ToAccessTokenScopes(token.Scopes),
			IssuedAt:  types.ToProtoTimestamp(token.IssuedAt),
			ExpiresAt: types.ToProtoTimestamp(token.ExpiresAt),
			RevokedAt: types.ToProtoTimestamp(token.RevokedAt),
		})
	}

	return res, nil
}

func (s *AdminServer) RevokeAccessToken(ctx context.Context, req *adminv1.RevokeAccessTokenRequest) (*adminv1.RevokeAccessTokenResponse, error) {
	tokenID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid access token id")
	}

	err = s.dbQuerier.RevokeAccessToken(ctx, s.pgxPool, tokenID)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke access tokens")
		return nil, status.Error(codes.Internal, "failed to revoke access tokens")
	}

	return &adminv1.RevokeAccessTokenResponse{}, nil
}
