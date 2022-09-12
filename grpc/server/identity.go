package server

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	identityv1 "github.com/nanzhong/tstr/api/identity/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityServer struct {
	identityv1.UnimplementedIdentityServiceServer

	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
}

func NewIdentityServer(pgxPool *pgxpool.Pool) identityv1.IdentityServiceServer {
	return &IdentityServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
	}
}

func (s *IdentityServer) Identity(ctx context.Context, r *identityv1.IdentityRequest) (*identityv1.IdentityResponse, error) {
	_, tokenHash, err := auth.AccessTokenFromContext(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get identity from metadata")
		return nil, status.Error(codes.Internal, "failed to get identity")
	}

	auth, err := db.New().AuthAccessToken(ctx, s.pgxPool, tokenHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to get identity")
		return nil, status.Error(codes.Internal, "failed to get identity")
	}

	return &identityv1.IdentityResponse{
		Scopes:             types.ToAccessTokenScopes(auth.Scopes),
		NamespaceSelectors: auth.NamespaceSelectors,
		IssuedAt:           types.ToProtoTimestamp(auth.IssuedAt),
		ExpiresAt:          types.ToProtoTimestamp(auth.ExpiresAt),
	}, nil
}
