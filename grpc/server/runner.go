package server

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerServer struct {
	dbQuerier db.Querier

	runner.UnimplementedRunnerServiceServer
}

func NewRunnerServer(dbQuerier db.Querier) runner.RunnerServiceServer {
	return &RunnerServer{
		dbQuerier: dbQuerier,
	}
}

func (s *RunnerServer) RegisterRunner(ctx context.Context, req *runner.RegisterRunnerRequest) (*runner.RegisterRunnerResponse, error) {
	var (
		accept pgtype.JSONB
		reject pgtype.JSONB
	)

	if err := accept.Set(req.AcceptTestLabelSelectors); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid accept test label selectors")
	}
	if err := reject.Set(req.RejectTestLabelSelectors); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reject test label selectors")
	}

	regRunner, err := s.dbQuerier.RegisterRunner(ctx, db.RegisterRunnerParams{
		Name:                     req.Name,
		AcceptTestLabelSelectors: accept,
		RejectTestLabelSelectors: reject,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to register runner in db")
		return nil, status.Error(codes.Internal, "failed to register runner")
	}

	return &runner.RegisterRunnerResponse{
		Runner: &common.Runner{
			Id:                       regRunner.ID,
			Name:                     regRunner.Name,
			AcceptTestLabelSelectors: map[string]string{},
			RejectTestLabelSelectors: map[string]string{},
			RegisteredAt:             toProtoTimestamp(regRunner.RegisteredAt),
			LastHeartbeatAt:          toProtoTimestamp(regRunner.LastHeartbeatAt),
		},
	}, nil
}

func (s *RunnerServer) NextRun(ctx context.Context, req *runner.NextRunRequest) (*runner.NextRunResponse, error) {
	return nil, nil
}

func (s *RunnerServer) SubmitRun(server runner.RunnerService_SubmitRunServer) error {
	return nil
}
