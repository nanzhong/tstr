package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestRunnerServer(t *testing.T) (*RunnerServer, *db.MockQuerier){
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &RunnerServer{
		dbQuerier: mockQuerier,
	}, mockQuerier
}


func TestRunnerServer_ResgisterRunner(t *testing.T) {
	/*tests := []struct {
		name				string
		responseCode		codes.Code
		errMsg				string
		mockQuerierReturn	func(*commonv1.Runner) (db.RegisterRunnerParams, error)

	}{
		{
		name: "",
		responseCode: code.OK,
		errMsg: "",
		mockQuerierReturn: func(*commonv1.Runner) (db.RegisterRunnerParams, error){
			token, _ := s.dbQuerier.AuthAccessToken(ctx, s.pgxPool, tokenHash)

			return commonv1.Runner{
				Id:                       regRunner.ID.String(),
				Name:                     regRunner.Name,
				NamespaceSelectors:       regRunner.NamespaceSelectors,
				AcceptTestLabelSelectors: acceptSelectors,
				RejectTestLabelSelectors: rejectSelectors,
				RegisteredAt:             types.ToProtoTimestamp(regRunner.RegisteredAt),
				LastHeartbeatAt:          types.ToProtoTimestamp(regRunner.LastHeartbeatAt),

			} 
		}
		*/

	server, mockQuerier := newTestRunnerServer(t)
	tokenString := "token"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(auth.MDKeyAuth, "bearer "+tokenString))
	token := newAccessTokenBuilder().withRevokedAt().build()
	runnerName := "runner"
	runnerId, _ := uuid.Parse(runnerName) 
	tokenID, _ := uuid.Parse(token.Id)
	tokenHash := auth.HashToken(tokenString)

	var (
		accept pgtype.JSONB
		reject pgtype.JSONB
		acceptSelectors pgtype.JSONB
		rejectSelectors pgtype.JSONB
	)

	mockQuerier.EXPECT().AuthAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenHash).Return(db.AuthAccessTokenRow{
		ID:                 tokenID,
		Name:               token.Name,
		NamespaceSelectors: token.NamespaceSelectors,
		Scopes:             []string{"admin"},
		IssuedAt:           types.FromProtoTimestampAsNullTime(token.IssuedAt),
		ExpiresAt:          types.FromProtoTimestampAsNullTime(token.ExpiresAt),
	}, nil)

	mockQuerier.EXPECT().RegisterRunner(ctx, gomock.AssignableToTypeOf(server.pgxPool), db.RegisterRunnerParams{
		Name:                     runnerName,
		NamespaceSelectors:       token.NamespaceSelectors,
		AcceptTestLabelSelectors: accept,
		RejectTestLabelSelectors: reject,
	}).Return(db.Runner{
		ID:                       runnerId,
		Name:                     runnerName,
		NamespaceSelectors:       token.NamespaceSelectors,
		AcceptTestLabelSelectors: acceptSelectors,
		RejectTestLabelSelectors: rejectSelectors,
		RegisteredAt:             types.FromProtoTimestampAsNullTime(token.IssuedAt),
		LastHeartbeatAt:          types.FromProtoTimestampAsNullTime(token.IssuedAt),
	}, nil)

	request := &runnerv1.RegisterRunnerRequest{Name: "TestRunner"}
	res , err := server.RegisterRunner(ctx, request)
	require.NoError(t, err)
	fmt.Println(res)
}