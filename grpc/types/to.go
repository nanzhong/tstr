package types

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUUIDString(uuid pgtype.UUID) string {
	if uuid.Status == pgtype.Null {
		return ""
	}
	var uuidString string
	_ = uuid.AssignTo(&uuidString)
	return uuidString
}

func ToAccessTokenScope(scope db.AccessTokenScope) common.AccessToken_Scope {
	switch scope {
	case db.AccessTokenScopeAdmin:
		return common.AccessToken_ADMIN
	case db.AccessTokenScopeControlR:
		return common.AccessToken_CONTROL_R
	case db.AccessTokenScopeControlRw:
		return common.AccessToken_CONTROL_RW
	case db.AccessTokenScopeRunner:
		return common.AccessToken_RUNNER
	default:
		return common.AccessToken_UNKNOWN
	}
}

func ToAccessTokenScopes(scopes []db.AccessTokenScope) []common.AccessToken_Scope {
	var protoScopes []common.AccessToken_Scope
	for _, s := range scopes {
		protoScopes = append(protoScopes, ToAccessTokenScope(s))
	}
	return protoScopes
}

func ToProtoTimestamp(ts interface{}) *timestamppb.Timestamp {
	switch tt := ts.(type) {
	case time.Time:
		return timestamppb.New(tt)
	case sql.NullTime:
		if tt.Valid {
			return timestamppb.New(tt.Time)
		}
		return nil
	default:
		panic("unexpected type for time")
	}
}
