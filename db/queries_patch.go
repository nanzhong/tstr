package db

import "github.com/jackc/pgtype"

// XXX there is bug in pggen for array enum types where this doesn't seem to be generated

func newRunResultEnum() pgtype.ValueTranscoder {
	return pgtype.NewEnumType("run_result", []string{string(RunResultUnknown), string(RunResultPass), string(RunResultFail), string(RunResultError)})
}

func newAccessTokenScopeEnum() pgtype.ValueTranscoder {
	return pgtype.NewEnumType("access_token_scope", []string{string(AccessTokenScopeAdmin), string(AccessTokenScopeControlR), string(AccessTokenScopeControlRW)})
}
