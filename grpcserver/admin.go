package grpcserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/api/admin/v1"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AdminServer struct {
	admin.UnimplementedAdminServiceServer

	dbQuerier db.Querier
}

func NewAdminServer(dbQuerier db.Querier) admin.AdminServiceServer {
	return &AdminServer{
		dbQuerier: dbQuerier,
	}
}

func (s *AdminServer) IssueAccessToken(ctx context.Context, req *admin.IssueAccessTokenRequest) (*admin.IssueAccessTokenResponse, error) {
	var scopes []db.AccessTokenScope
	for _, s := range req.Scopes {
		switch s {
		case common.AccessToken_ADMIN:
			scopes = append(scopes, db.AccessTokenScopeAdmin)
		case common.AccessToken_CONTROL_R:
			scopes = append(scopes, db.AccessTokenScopeControlR)
		case common.AccessToken_CONTROL_RW:
			scopes = append(scopes, db.AccessTokenScopeControlRW)
		default:
			log.Error().Int32("scope", int32(s)).Msg("invalid scope")
			return nil, status.Error(codes.InvalidArgument, "failed to issue access token")
		}
	}

	var expiresAt pgtype.Timestamptz
	if req.ValidDuration != nil {
		expiresAt.Set(time.Now().UTC().Add(req.ValidDuration.AsDuration()))
	}
	tokenData := make([]byte, 32)
	_, err := rand.Read(tokenData)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token value")
		return nil, status.Error(codes.Internal, "failed to issue access token")
	}
	token := hex.EncodeToString(tokenData)

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash token")
		return nil, status.Error(codes.Internal, "failed to issue access token")
	}

	issuedToken, err := s.dbQuerier.IssueAccessToken(ctx, db.IssueAccessTokenParams{
		Name:      req.Name,
		TokenHash: string(tokenHash),
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create access token in db")
		return nil, status.Error(codes.Internal, "failed to issue access token")
	}

	return &admin.IssueAccessTokenResponse{
		AccessToken: &common.AccessToken{
			Id:        issuedToken.ID,
			Name:      issuedToken.Name,
			Scopes:    toProtoScopes(issuedToken.Scopes),
			Token:     token,
			IssuedAt:  toProtoTimestamp(issuedToken.IssuedAt),
			ExpiresAt: toProtoTimestamp(issuedToken.ExpiresAt),
		},
	}, nil
}

func (s *AdminServer) GetAccessToken(ctx context.Context, req *admin.GetAccessTokenRequest) (*admin.GetAccessTokenResponse, error) {
	token, err := s.dbQuerier.GetAccessToken(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed to query for access token")
		return nil, status.Error(codes.Internal, "failed to get access token")
	}

	return &admin.GetAccessTokenResponse{
		AccessToken: &common.AccessToken{
			Id:        token.ID,
			Name:      token.Name,
			Scopes:    toProtoScopes(token.Scopes),
			IssuedAt:  toProtoTimestamp(token.IssuedAt),
			ExpiresAt: toProtoTimestamp(token.ExpiresAt),
			RevokedAt: toProtoTimestamp(token.RevokedAt),
		},
	}, nil
}

func (s *AdminServer) ListAccessTokens(ctx context.Context, req *admin.ListAccessTokensRequest) (*admin.ListAccessTokensResponse, error) {
	tokens, err := s.dbQuerier.ListAccessTokens(ctx, req.IncludeExpired, req.IncludeRevoked)
	if err != nil {
		log.Error().Err(err).Msg("failed to query for access tokens")
		return nil, status.Error(codes.Internal, "failed to list access tokens")
	}

	res := &admin.ListAccessTokensResponse{}
	for _, token := range tokens {
		res.AccessTokens = append(res.AccessTokens, &common.AccessToken{
			Id:        token.ID,
			Name:      token.Name,
			Scopes:    toProtoScopes(token.Scopes),
			IssuedAt:  toProtoTimestamp(token.IssuedAt),
			ExpiresAt: toProtoTimestamp(token.ExpiresAt),
			RevokedAt: toProtoTimestamp(token.RevokedAt),
		})
	}

	return res, nil
}

func (s *AdminServer) RevokeAccessToken(ctx context.Context, req *admin.RevokeAccessTokenRequest) (*admin.RevokeAccessTokenResponse, error) {
	_, err := s.dbQuerier.RevokeAccessToken(ctx, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed to revoke access tokens")
		return nil, status.Error(codes.Internal, "failed to revoke access tokens")
	}

	return &admin.RevokeAccessTokenResponse{}, nil
}

func (s *AdminServer) ApproveRunner(ctx context.Context, req *admin.ApproveRunnerRequest) (*admin.ApproveRunnerResponse, error) {
	return nil, nil
}

func (s *AdminServer) DenyRunner(ctx context.Context, req *admin.DenyRunnerRequest) (*admin.DenyRunnerResponse, error) {
	return nil, nil
}

func toProtoScopes(scopes []db.AccessTokenScope) []common.AccessToken_Scope {
	var protoScopes []common.AccessToken_Scope
	for _, s := range scopes {
		switch s {
		case db.AccessTokenScopeAdmin:
			protoScopes = append(protoScopes, common.AccessToken_ADMIN)
		case db.AccessTokenScopeControlR:
			protoScopes = append(protoScopes, common.AccessToken_CONTROL_R)
		case db.AccessTokenScopeControlRW:
			protoScopes = append(protoScopes, common.AccessToken_CONTROL_RW)
		}
	}
	return protoScopes
}

func toProtoTimestamp(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if ts.Status == pgtype.Null {
		return nil
	}

	var t time.Time
	if err := ts.AssignTo(&t); err != nil {
		// This should never fail...
		log.Error().Err(err).Msg("converting pgtype.timestamptz to time.Time")
		return nil
	}

	return timestamppb.New(t)
}
