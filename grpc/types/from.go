package types

import (
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
)

func FromAccessTokenScope(scope common.AccessToken_Scope) db.AccessTokenScope {
	switch scope {
	case common.AccessToken_ADMIN:
		return db.AccessTokenScopeAdmin
	case common.AccessToken_CONTROL_R:
		return db.AccessTokenScopeControlR
	case common.AccessToken_CONTROL_RW:
		return db.AccessTokenScopeControlRw
	case common.AccessToken_RUNNER:
		return db.AccessTokenScopeRunner
	default:
		panic("unknown access token scope:" + scope.String())
	}
}

func FromAccessTokenScopes(scopes []common.AccessToken_Scope) []db.AccessTokenScope {
	var dbScopes []db.AccessTokenScope
	for _, s := range scopes {
		dbScopes = append(dbScopes, FromAccessTokenScope(s))
	}
	return dbScopes
}
