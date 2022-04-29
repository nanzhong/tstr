package grpcserver

import "github.com/nanzhong/tstr/api/runner/v1"

type RunnerServer struct {
	runner.UnimplementedRunnerServiceServer
}

func NewRunnerServer() runner.RunnerServiceServer {
	return &RunnerServer{}
}
