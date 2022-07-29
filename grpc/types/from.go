package types

import (
	"database/sql"

	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromAccessTokenScope(scope commonv1.AccessToken_Scope) db.AccessTokenScope {
	switch scope {
	case commonv1.AccessToken_ADMIN:
		return db.AccessTokenScopeAdmin
	case commonv1.AccessToken_CONTROL_R:
		return db.AccessTokenScopeControlR
	case commonv1.AccessToken_CONTROL_RW:
		return db.AccessTokenScopeControlRw
	case commonv1.AccessToken_RUNNER:
		return db.AccessTokenScopeRunner
	case commonv1.AccessToken_DATA:
		return db.AccessTokenScopeData
	default:
		panic("unknown access token scope:" + scope.String())
	}
}

func FromAccessTokenScopes(scopes []commonv1.AccessToken_Scope) []db.AccessTokenScope {
	var dbScopes []db.AccessTokenScope
	for _, s := range scopes {
		dbScopes = append(dbScopes, FromAccessTokenScope(s))
	}
	return dbScopes
}

func FromRunResult(result commonv1.Run_Result) db.RunResult {
	switch result {
	case commonv1.Run_ERROR:
		return db.RunResultError
	case commonv1.Run_FAIL:
		return db.RunResultFail
	case commonv1.Run_PASS:
		return db.RunResultPass
	default:
		return db.RunResultUnknown
	}
}

// FromRunResults converts pb run results to db run results.
// XXX There is currently a bug preventing enmu arrays from being usable. Using
// []string until this is resolved.
// https://github.com/kyleconroy/sqlc/issues/1256
func FromRunResults(results []commonv1.Run_Result) []string {
	var dbResults []string
	for _, r := range results {
		dbResults = append(dbResults, string(FromRunResult(r)))
	}
	return dbResults
}

func FromProtoTimestampAsNullTime(ts *timestamppb.Timestamp) sql.NullTime {
	if ts == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true, Time: ts.AsTime()}
}

func FromProtoTestRunConfig(rc *commonv1.Test_RunConfig) db.TestRunConfig {
	//	env := pgtype.JSONB{}
	//	if err := env.Set(rc.Env); err != nil {
	//		return nil, fmt.Errorf("parsing env: %w", err)
	//	}

	return db.TestRunConfig{
		ContainerImage: rc.ContainerImage,
		Command:        rc.Command,
		Args:           rc.Args,
		Env:            rc.Env,
		TimeoutSeconds: uint(rc.Timeout.Seconds),
	}
}
