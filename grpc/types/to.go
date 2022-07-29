package types

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"google.golang.org/protobuf/types/known/durationpb"
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

// ToAccessTokenScope converts db access token scope to the pb type.
// XXX There is currently a bug preventing enmu arrays from being usable. Using
// []string until this is resolved.
// https://github.com/kyleconroy/sqlc/issues/1256
func ToAccessTokenScope(scope string) commonv1.AccessToken_Scope {
	switch scope {
	case string(db.AccessTokenScopeAdmin):
		return commonv1.AccessToken_ADMIN
	case string(db.AccessTokenScopeControlR):
		return commonv1.AccessToken_CONTROL_R
	case string(db.AccessTokenScopeControlRw):
		return commonv1.AccessToken_CONTROL_RW
	case string(db.AccessTokenScopeRunner):
		return commonv1.AccessToken_RUNNER
	case string(db.AccessTokenScopeData):
		return commonv1.AccessToken_DATA
	default:
		return commonv1.AccessToken_UNKNOWN
	}
}

// ToAccessTokenScopes converts db access token scopes to the pb type.
// XXX There is currently a bug preventing enmu arrays from being usable. Using
// []string until this is resolved.
// https://github.com/kyleconroy/sqlc/issues/1256
func ToAccessTokenScopes(scopes []string) []commonv1.AccessToken_Scope {
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

func ToProtoResultData(data pgtype.JSONB) map[string]string {
	r := map[string]string{}
	data.AssignTo(&r)
	return r
}

func ToProtoTestRunConfig(rc db.TestRunConfig) *commonv1.Test_RunConfig {
	return &commonv1.Test_RunConfig{
		ContainerImage: rc.ContainerImage,
		Command:        rc.Command,
		Args:           rc.Args,
		Env:            rc.Env,
		Timeout:        durationpb.New(time.Duration(rc.TimeoutSeconds) * time.Second),
	}
}
