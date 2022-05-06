package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const mdAuthKey = "authorization"

var scopeAuthorizations = map[string][]common.AccessToken_Scope{
	"/tstr.control.v1.ControlService/RegisterTest":     {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/UpdateTest":       {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/GetTest":          {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ListTests":        {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ArchiveTest":      {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/DefineTestSuite":  {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/UpdateSuite":      {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/GetTestSuite":     {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ListTestSuites":   {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ArchiveTestSuite": {common.AccessToken_CONTROL_RW},
	"/tstr.control.v1.ControlService/GetRun":           {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ListRuns":         {common.AccessToken_CONTROL_RW, common.AccessToken_CONTROL_R},
	"/tstr.control.v1.ControlService/ScheduleRun":      {common.AccessToken_CONTROL_RW},

	"/tstr.admin.v1.AdminService/IssueAccessToken":  {common.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/GetAccessToken":    {common.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/ListAccessTokens":  {common.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/RevokeAccessToken": {common.AccessToken_ADMIN},

	"/tstr.runner.v1.RunnerService/RegisterRunner": {common.AccessToken_RUNNER},
	"/tstr.runner.v1.RunnerService/NextRun":        {common.AccessToken_RUNNER},
	"/tstr.runner.v1.RunnerService/SubmitRun":      {common.AccessToken_RUNNER},
}

// TODO We shouldn't reach out to the db each time to auth, especially when
// these results are easy to cache/invalidate.

func UnaryServerInterceptor(dbQuerier db.Querier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, exists := metadata.FromIncomingContext(ctx)
		if !exists {
			return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: missing access token")
		}

		_, tokenHash, err := tokenFromMD(md)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		validScopes := scopeAuthorizations[info.FullMethod]
		fmt.Println(tokenHash)
		allowed, err := dbQuerier.ValidateAccessToken(ctx, tokenHash, toDBScopes(validScopes))
		if err != nil {
			log.Error().Err(err).Msg("failed to validate access token")
			return nil, status.Error(codes.Internal, "failed to authenticate request")
		}

		if !allowed {
			return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(dbQuerier db.Querier) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, exists := metadata.FromIncomingContext(ss.Context())
		if !exists {
			return status.Error(codes.Unauthenticated, "failed to authenticate request: missing access token")
		}

		_, tokenHash, err := tokenFromMD(md)
		if err != nil {
			return status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		validScopes := scopeAuthorizations[info.FullMethod]
		allowed, err := dbQuerier.ValidateAccessToken(ss.Context(), tokenHash, toDBScopes(validScopes))
		if err != nil {
			log.Error().Err(err).Msg("failed to validate access token")
			return status.Error(codes.Internal, "failed to authenticate request")
		}

		if !allowed {
			return status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		return handler(srv, ss)
	}
}

func UnaryClientInterceptor(accessToken string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, mdAuthKey, "bearer "+accessToken)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(accessToken string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, mdAuthKey, "bearer "+accessToken)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func toDBScope(scope common.AccessToken_Scope) db.AccessTokenScope {
	switch scope {
	case common.AccessToken_ADMIN:
		return db.AccessTokenScopeAdmin
	case common.AccessToken_CONTROL_R:
		return db.AccessTokenScopeControlR
	case common.AccessToken_CONTROL_RW:
		return db.AccessTokenScopeControlRW
	case common.AccessToken_RUNNER:
		return db.AccessTokenScopeRunner
	case common.AccessToken_UNKNOWN:
		// This should never happen and is an indication that an endpoint is not
		// configured in the scope authorization map.
		panic("endpoint not scoped")
	default:
		// This should never happen and is an indication that proto scopes and db
		// scopes are not in sync.
		panic("missing scope definition")
	}
}

func toDBScopes(scopes []common.AccessToken_Scope) []db.AccessTokenScope {
	dbScopes := make([]db.AccessTokenScope, len(scopes))
	for i, s := range scopes {
		dbScopes[i] = toDBScope(s)
	}
	return dbScopes
}

func tokenFromMD(md metadata.MD) (string, string, error) {
	vals := md.Get(mdAuthKey)
	if vals == nil || len(vals) != 1 {
		return "", "", errors.New("invalid access token")
	}
	parts := strings.Split(strings.ToLower(vals[0]), " ")
	if len(parts) != 2 || parts[0] != "bearer" {
		return "", "", errors.New("invalid access token")
	}

	token := parts[1]
	tokenHashBytes := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(tokenHashBytes[:])

	return token, tokenHash, nil
}
