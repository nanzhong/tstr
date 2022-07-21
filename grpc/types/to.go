package types

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
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

func ToAccessTokenScope(scope db.AccessTokenScope) commonv1.AccessToken_Scope {
	switch scope {
	case db.AccessTokenScopeAdmin:
		return commonv1.AccessToken_ADMIN
	case db.AccessTokenScopeControlR:
		return commonv1.AccessToken_CONTROL_R
	case db.AccessTokenScopeControlRw:
		return commonv1.AccessToken_CONTROL_RW
	case db.AccessTokenScopeRunner:
		return commonv1.AccessToken_RUNNER
	case db.AccessTokenScopeData:
		return commonv1.AccessToken_DATA
	default:
		return commonv1.AccessToken_UNKNOWN
	}
}

func ToAccessTokenScopes(scopes []db.AccessTokenScope) []commonv1.AccessToken_Scope {
	var protoScopes []commonv1.AccessToken_Scope
	for _, s := range scopes {
		protoScopes = append(protoScopes, ToAccessTokenScope(s))
	}
	return protoScopes
}

func ToRunResult(result db.RunResult) commonv1.Run_Result {
	switch result {
	case db.RunResultError:
		return commonv1.Run_ERROR
	case db.RunResultFail:
		return commonv1.Run_FAIL
	case db.RunResultPass:
		return commonv1.Run_PASS
	default:
		return commonv1.Run_UNKNOWN
	}
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

func ToRunLogs(logs []db.RunLog) []*commonv1.Run_Log {
	var pbLogs []*commonv1.Run_Log
	for _, l := range logs {
		pbLogs = append(pbLogs, &commonv1.Run_Log{
			Time:       l.Time,
			OutputType: commonv1.Run_Log_Output(commonv1.Run_Log_Output_value[l.Type]),
			Data:       l.Data,
		})
	}
	return pbLogs
}
