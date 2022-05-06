package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/nanzhong/tstr/api/runner/v1"
	"github.com/rs/zerolog/log"
)

var errRunnerRevoked = errors.New("runner revoked")

type Runner struct {
	client runner.RunnerServiceClient

	id                   string
	name                 string
	allowLabelSelectors  map[string]string
	rejectLabelSelectors map[string]string

	doneCh       chan struct{}
	stopCh       chan struct{}
	stopCancelFn context.CancelFunc
}

func New(
	client runner.RunnerServiceClient,
	name string,
	allowLabelSelectors map[string]string,
	rejectLabelSelectors map[string]string,
) *Runner {
	return &Runner{
		client: client,

		name:                 name,
		allowLabelSelectors:  allowLabelSelectors,
		rejectLabelSelectors: rejectLabelSelectors,

		doneCh: make(chan struct{}),
		stopCh: make(chan struct{}),
	}
}

func (r *Runner) Run() error {
	defer close(r.doneCh)

	var ctx context.Context
	ctx, r.stopCancelFn = context.WithCancel(context.Background())
	regRes, err := r.client.RegisterRunner(ctx, &runner.RegisterRunnerRequest{
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
