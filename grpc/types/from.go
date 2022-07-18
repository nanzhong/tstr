package types

import (
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
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
