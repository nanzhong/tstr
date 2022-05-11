package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/client"
	"github.com/nanzhong/tstr/api/runner/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

var errRunnerRevoked = errors.New("runner revoked")

type Runner struct {
	runnerClient runner.RunnerServiceClient
	dockerClient *client.Client
	retryDelay   time.Duration

	id                   string
	name                 string
	allowLabelSelectors  map[string]string
	rejectLabelSelectors map[string]string

	doneCh       chan struct{}
	stopCh       chan struct{}
	stopCancelFn context.CancelFunc
}

func New(
	runnerClient runner.RunnerServiceClient,
	name string,
	allowLabelSelectors map[string]string,
	rejectLabelSelectors map[string]string,
) (*Runner, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to configure docker client: %w", err)
	}

	return &Runner{
		runnerClient: runnerClient,
		dockerClient: dockerClient,
		retryDelay:   5 * time.Second,

		name:                 name,
		allowLabelSelectors:  allowLabelSelectors,
		rejectLabelSelectors: rejectLabelSelectors,

		doneCh: make(chan struct{}),
		stopCh: make(chan struct{}),
	}, nil
}

func (r *Runner) Run() error {
	defer close(r.doneCh)

	var ctx context.Context
	ctx, r.stopCancelFn = context.WithCancel(context.Background())

	_, err := r.dockerClient.Ping(ctx)
	if err != nil {
		return fmt.Errorf("docker not available, unable to ping: %w", err)
	}

	regRes, err := r.runnerClient.RegisterRunner(ctx, &runner.RegisterRunnerRequest{
		Name:                     r.name,
		AcceptTestLabelSelectors: r.allowLabelSelectors,
		RejectTestLabelSelectors: r.rejectLabelSelectors,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to register runner")
		return fmt.Errorf("registering runner: %w", err)
	}

	r.id = regRes.Runner.Id

	for {
		select {
		case <-r.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		default:
			runRes, err := r.runnerClient.NextRun(ctx, &runner.NextRunRequest{
				Id: r.id,
			})
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
					log.Info().Msg("no pending runs, trying again later")
				} else {
					log.Error().Err(err).Msg("failed to get the next run")
				}
				time.Sleep(r.retryDelay)
			}

			fmt.Println(protojson.Format(runRes))
		}
	}
}

func (r *Runner) Stop(ctx context.Context) {
	close(r.stopCh)

	select {
	case <-r.doneCh:
	case <-ctx.Done():
		if r.stopCancelFn != nil {
			r.stopCancelFn()
		}
	}
}
