package server

import (
	"context"
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
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

	token, err := s.dbQuerier.AuthAccessToken(ctx, s.pgxPool, tokenHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to get identity")
		return nil, status.Error(codes.Internal, "failed to get identity")
	}

	namespaces, err := s.dbQuerier.ListAllNamespaces(ctx, s.pgxPool)
	if err != nil {
		log.Error().Err(err).Msg("failed to list all namespaces")
		return nil, status.Error(codes.Internal, "failed to get identity")
	}

	var nsREs []*regexp.Regexp
	for _, nsSel := range token.NamespaceSelectors {
		re, err := regexp.Compile(nsSel)
		if err != nil {
			log.Error().
				Err(err).
				Stringer("access_token_id", token.ID).
				Str("namespace_selector", nsSel).
				Msg("failed to compile namespace selector")
			return nil, status.Error(codes.Internal, "failed to get identity")
		}

		nsREs = append(nsREs, re)
	}

	var accessibleNamespaces []string
	for _, ns := range namespaces {
		for _, nsRE := range nsREs {
			if nsRE.MatchString(ns) {
				accessibleNamespaces = append(accessibleNamespaces, ns)
			}
		}
	}

	return &identityv1.IdentityResponse{
		AccessToken: &commonv1.AccessToken{
			Id:                 token.ID.String(),
			Name:               token.Name,
			NamespaceSelectors: token.NamespaceSelectors,
			Scopes:             types.ToAccessTokenScopes(token.Scopes),
			IssuedAt:           types.ToProtoTimestamp(token.IssuedAt),
			ExpiresAt:          types.ToProtoTimestamp(token.ExpiresAt),
		},
		AccessibleNamespaces: accessibleNamespaces,
	}, nil
}
