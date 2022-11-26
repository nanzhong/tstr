package auth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	MDKeyAuth      = "authorization"
	MDKeyNamespace = "namespace"
)

// TODO We can probably be a bit smarter/less verbose here and instead of direct
// string matches on the full method, build up the set of allowable tokens given
// a full method and a list of regexes (or something like that)
var scopeAuthorizations = map[string][]commonv1.AccessToken_Scope{
	"/tstr.identity.v1.IdentityService/Identity": {
		commonv1.AccessToken_ADMIN,
		commonv1.AccessToken_CONTROL,
		commonv1.AccessToken_RUNNER,
		commonv1.AccessToken_DATA,
	},

	"/tstr.control.v1.ControlService/RegisterTest": {commonv1.AccessToken_CONTROL},
	"/tstr.control.v1.ControlService/UpdateTest":   {commonv1.AccessToken_CONTROL},
	"/tstr.control.v1.ControlService/DeleteTest":   {commonv1.AccessToken_CONTROL},
	"/tstr.control.v1.ControlService/ScheduleRun":  {commonv1.AccessToken_CONTROL},

	"/tstr.admin.v1.AdminService/IssueAccessToken":  {commonv1.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/GetAccessToken":    {commonv1.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/ListAccessTokens":  {commonv1.AccessToken_ADMIN},
	"/tstr.admin.v1.AdminService/RevokeAccessToken": {commonv1.AccessToken_ADMIN},

	"/tstr.runner.v1.RunnerService/RegisterRunner": {commonv1.AccessToken_RUNNER},
	"/tstr.runner.v1.RunnerService/NextRun":        {commonv1.AccessToken_RUNNER},
	"/tstr.runner.v1.RunnerService/SubmitRun":      {commonv1.AccessToken_RUNNER},

	"/tstr.data.v1.DataService/GetTest":         {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/QueryTests":      {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/GetTestSuite":    {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/QueryTestSuites": {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/GetRun":          {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/QueryRuns":       {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/SummarizeRuns":   {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/GetRunner":       {commonv1.AccessToken_DATA},
	"/tstr.data.v1.DataService/QueryRunners":    {commonv1.AccessToken_DATA},
}

// TODO We shouldn't reach out to the db each time to auth, especially when
// these results are easy to cache/invalidate.

func UnaryServerInterceptor(pgxPool *pgxpool.Pool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, exists := metadata.FromIncomingContext(ctx)
		if !exists {
			return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: missing access token")
		}

		_, tokenHash, err := AccessTokenFromMD(md)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		validScopes := scopeAuthorizations[info.FullMethod]
		auth, err := db.New().AuthAccessToken(ctx, pgxPool, tokenHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
			}

			log.Error().Err(err).Msg("failed to validate access token")
			return nil, status.Error(codes.Internal, "failed to authenticate request")
		}

		if !authorizeScope(auth, validScopes) {
			return nil, status.Error(codes.PermissionDenied, "failed to authenticate request: invalid access token scopes")
		}

		if strings.HasPrefix(info.FullMethod, "/tstr.control.v1") ||
			strings.HasPrefix(info.FullMethod, "/tstr.data.v1") {
			namespace, err := NamespaceFromMD(md)
			if err != nil {
				return nil, status.Error(codes.PermissionDenied, "failed to authorize request: invalid namespace")
			}

			allowed, err := authorizeNamespace(auth, namespace)
			if err != nil {
				return nil, status.Error(codes.Internal, "failed to authorize request")
			}
			if !allowed {
				return nil, status.Error(codes.PermissionDenied, "failed to authorize request: invalid namespace")
			}
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(pgxPool *pgxpool.Pool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, exists := metadata.FromIncomingContext(ss.Context())
		if !exists {
			return status.Error(codes.Unauthenticated, "failed to authenticate request: missing access token")
		}

		_, tokenHash, err := AccessTokenFromMD(md)
		if err != nil {
			return status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
		}

		validScopes := scopeAuthorizations[info.FullMethod]
		auth, err := db.New().AuthAccessToken(ss.Context(), pgxPool, tokenHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return status.Error(codes.Unauthenticated, "failed to authenticate request: invalid access token")
			}

			log.Error().Err(err).Msg("failed to validate access token")
			return status.Error(codes.Internal, "failed to authenticate request")
		}

		if !authorizeScope(auth, validScopes) {
			return status.Error(codes.PermissionDenied, "failed to authenticate request: invalid access token scopes")
		}

		if strings.HasPrefix(info.FullMethod, "/tstr.control.v1") ||
			strings.HasPrefix(info.FullMethod, "/tstr.data.v1") {
			namespace, err := NamespaceFromMD(md)
			if err != nil {
				return status.Error(codes.PermissionDenied, "failed to authorize request: invalid namespace")
			}

			allowed, err := authorizeNamespace(auth, namespace)
			if err != nil {
				return status.Error(codes.Internal, "failed to authorize request")
			}
			if !allowed {
				return status.Error(codes.PermissionDenied, "failed to authorize request: invalid namespace")
			}
		}

		return handler(srv, ss)
	}
}

func authorizeScope(auth db.AuthAccessTokenRow, validScopes []commonv1.AccessToken_Scope) bool {
	for _, vs := range types.FromAccessTokenScopes(validScopes) {
		for _, s := range auth.Scopes {
			if string(vs) == s {
				return true
			}
		}
	}
	return false
}

func authorizeNamespace(auth db.AuthAccessTokenRow, namespace string) (bool, error) {
	namespaceAllowed := false
	for _, nsSel := range auth.NamespaceSelectors {
		re, err := regexp.Compile(nsSel)
		if err != nil {
			return false, status.Error(codes.Internal, "failed to authorize request: error validating namespace")
		}
		if re.MatchString(namespace) {
			namespaceAllowed = true
			break
		}
	}
	if !namespaceAllowed {
		return false, status.Error(codes.PermissionDenied, "failed to authorize request: invalid namespace")
	}

	return true, nil
}

func UnaryClientInterceptor(accessToken string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, MDKeyAuth, "bearer "+accessToken)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamClientInterceptor(accessToken string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, MDKeyAuth, "bearer "+accessToken)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func AccessTokenFromMD(md metadata.MD) (string, string, error) {
	vals := md.Get(MDKeyAuth)
	if vals == nil || len(vals) != 1 {
		return "", "", errors.New("invalid access token")
	}
	parts := strings.Split(strings.ToLower(vals[0]), " ")
	if len(parts) != 2 || parts[0] != "bearer" {
		return "", "", errors.New("invalid access token")
	}

	token := parts[1]
	tokenHash := HashToken(token)

	return token, tokenHash, nil
}

func AccessTokenFromContext(ctx context.Context) (string, string, error) {
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return "", "", errors.New("context missing metadata")
	}

	return AccessTokenFromMD(md)
}

func NamespaceFromMD(md metadata.MD) (string, error) {
	vals := md.Get(MDKeyNamespace)
	if vals == nil || len(vals) != 1 {
		return "", errors.New("metadata missing namespace")
	}
	return vals[0], nil
}

func NamespaceFromContext(ctx context.Context) (string, error) {
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return "", errors.New("context missing metadata")
	}

	ns, err := NamespaceFromMD(md)
	if err != nil {
		return "", fmt.Errorf("getting namespace from context: %w", err)
	}
	return ns, nil
}

func HashToken(token string) string {
	tokenHashBytes := sha512.Sum512([]byte(token))
	return hex.EncodeToString(tokenHashBytes[:])
}
