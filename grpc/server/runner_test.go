package server

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/auth"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func newTestRunnerServer(t *testing.T) (*RunnerServer, *db.MockQuerier) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &RunnerServer{
		dbQuerier: mockQuerier,
	}, mockQuerier
}

func createRunnerData(namespaceSelectors []string) *commonv1.Runner {
	runnerName := "runner"
	runnerId, _ := uuid.Parse(runnerName)
	timestamp := types.ToProtoTimestamp(time.Now())
	return &commonv1.Runner{
		Id:                       runnerId.String(),
		Name:                     runnerName,
		NamespaceSelectors:       namespaceSelectors,
		AcceptTestLabelSelectors: map[string]string{"region": "nyc"},
		RejectTestLabelSelectors: map[string]string{"region": "syd"},
		RegisteredAt:             timestamp,
		LastHeartbeatAt:          timestamp,
	}
}

func convertTestLabelSelectors(testLabelSelectors map[string]string) pgtype.JSONB {
	var transformedLabel pgtype.JSONB
	transformedLabel.Set(testLabelSelectors)
	return transformedLabel
}

func TestRunnerServer_RegisterRunner(t *testing.T) {
	tests := []struct {
		name                    string
		responseCode            codes.Code
		errMsg                  string
		token                   *commonv1.AccessToken
		runner                  func([]string) *commonv1.Runner
		mockQuerierReturnAuth   func(*commonv1.AccessToken) (db.AuthAccessTokenRow, error)
		mockQuerierReturnRunner func(*commonv1.Runner) (db.Runner, error)
	}{
		{
			name:         "valid",
			responseCode: codes.OK,
			errMsg:       "",
			token:        newAccessTokenBuilder().build(),
			runner: func(namespaceSelectors []string) *commonv1.Runner {
				return createRunnerData(namespaceSelectors)
			},
			mockQuerierReturnAuth: func(token *commonv1.AccessToken) (db.AuthAccessTokenRow, error) {
				tokenID, _ := uuid.Parse(token.Id)
				return db.AuthAccessTokenRow{
					ID:                 tokenID,
					Name:               token.Name,
					NamespaceSelectors: token.NamespaceSelectors,
					Scopes:             []string{"admin"},
					IssuedAt:           types.FromProtoTimestampAsNullTime(token.IssuedAt),
					ExpiresAt:          types.FromProtoTimestampAsNullTime(token.ExpiresAt),
					RevokedAt:          types.FromProtoTimestampAsNullTime(token.RevokedAt),
				}, nil
			},
			mockQuerierReturnRunner: func(runner *commonv1.Runner) (db.Runner, error) {
				runnerID, _ := uuid.Parse(runner.Id)
				return db.Runner{
					ID:                       runnerID,
					Name:                     runner.Name,
					NamespaceSelectors:       runner.NamespaceSelectors,
					AcceptTestLabelSelectors: convertTestLabelSelectors(runner.AcceptTestLabelSelectors),
					RejectTestLabelSelectors: convertTestLabelSelectors(runner.RejectTestLabelSelectors),
					RegisteredAt:             types.FromProtoTimestampAsNullTime(runner.RegisteredAt),
					LastHeartbeatAt:          types.FromProtoTimestampAsNullTime(runner.LastHeartbeatAt),
				}, nil
			},
		},
	}

	server, mockQuerier := newTestRunnerServer(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(auth.MDKeyAuth, "bearer "+tt.token.Name))
			tokenHash := auth.HashToken(tt.token.Name)
			mockQuerier.EXPECT().AuthAccessToken(ctx, gomock.AssignableToTypeOf(server.pgxPool), tokenHash).Return(tt.mockQuerierReturnAuth(tt.token))

			runner := tt.runner(tt.token.NamespaceSelectors)
			mockQuerier.EXPECT().RegisterRunner(ctx, gomock.AssignableToTypeOf(server.pgxPool), db.RegisterRunnerParams{
				Name:                     runner.Name,
				NamespaceSelectors:       tt.token.NamespaceSelectors,
				AcceptTestLabelSelectors: convertTestLabelSelectors(runner.AcceptTestLabelSelectors),
				RejectTestLabelSelectors: convertTestLabelSelectors(runner.RejectTestLabelSelectors),
			}).Return(tt.mockQuerierReturnRunner(runner))

			request := &runnerv1.RegisterRunnerRequest{
				Name:                     runner.Name,
				AcceptTestLabelSelectors: runner.AcceptTestLabelSelectors,
				RejectTestLabelSelectors: runner.RejectTestLabelSelectors}
			res, err := server.RegisterRunner(ctx, request)

			assert.Equal(t, &runnerv1.RegisterRunnerResponse{Runner: runner}, res)
			require.NoError(t, err)
		})
	}
}
