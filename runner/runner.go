package runner

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errRunnerRevoked = errors.New("runner revoked")

var recordStart = []byte("\x1c\x1e")
var recordEnd = []byte("\x1e\x1c")

type Runner struct {
	runnerClient runnerv1.RunnerServiceClient
	dockerClient *client.Client
	retryDelay   time.Duration
	clock        clock.Clock

	id                   string
	name                 string
	allowLabelSelectors  map[string]string
	rejectLabelSelectors map[string]string

	doneCh       chan struct{}
	stopCh       chan struct{}
	stopCancelFn context.CancelFunc
}

func New(
	runnerClient runnerv1.RunnerServiceClient,
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
		clock:        clock.New(),

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

	regRes, err := r.runnerClient.RegisterRunner(ctx, &runnerv1.RegisterRunnerRequest{
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
		case <-r.clock.After(r.retryDelay):
			runRes, err := r.runnerClient.NextRun(ctx, &runnerv1.NextRunRequest{
				Id: r.id,
			})
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
					log.Info().Msg("no pending runs, trying again later")
				} else {
					log.Error().Err(err).Msg("failed to get the next run")
				}
				continue
			}

			if err := r.executeRun(ctx, runRes.Run); err != nil {
				log.Error().Err(err).Msg("failed to execute run")
			}
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

func (r Runner) executeRun(ctx context.Context, run *commonv1.Run) error {
	runLogger := log.With().
		Str("run_id", run.Id).
		Logger()

	runLogger.Info().Msg("starting run")

	stream, err := r.runnerClient.SubmitRun(ctx)
	if err != nil {
		return fmt.Errorf("creating submit run stream: %w", err)
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			runLogger.Error().Err(err).Msg("closing submit run stream")
		}
	}()

	if err := stream.Send(&runnerv1.SubmitRunRequest{
		Id:    r.id,
		RunId: run.Id,
	}); err != nil {
		return fmt.Errorf("registering run and runner id : %w", err)
	}
	runLogger.Debug().Msg("registered run and runner")

	runLogger.Info().
		Str("image", run.TestRunConfig.ContainerImage).
		Msg("pulling image")
	pullReader, err := r.dockerClient.ImagePull(ctx, run.TestRunConfig.ContainerImage, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pulling image (%s): %w", run.TestRunConfig.ContainerImage, err)
	}
	defer pullReader.Close()
	io.Copy(io.Discard, pullReader)

	var env []string
	for k, v := range run.TestRunConfig.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	env = append(env, "RECORD_START="+string(recordStart))
	env = append(env, "RECORD_END="+string(recordEnd))

	var cmd []string
	if run.TestRunConfig.Command != "" {
		cmd = append(cmd, run.TestRunConfig.Command)
	}
	for _, a := range run.TestRunConfig.Args {
		cmd = append(cmd, a)
	}

	containerName := fmt.Sprintf("run-%s", run.Id)
	runLogger.Info().Msg("creating container")
	createRes, err := r.dockerClient.ContainerCreate(ctx, &container.Config{
		Env:   env,
		Cmd:   cmd,
		Image: run.TestRunConfig.ContainerImage,
	}, nil, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}
	runLogger = runLogger.With().
		Str("container_id", createRes.ID).
		Str("container_name", containerName).
		Logger()

	defer func() {
		if err := r.dockerClient.ContainerRemove(ctx, createRes.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			runLogger.Error().Err(err).Msg("removing container for run")
		}
		runLogger.Info().Msg("cleaned up container")
	}()

	runLogger.Info().Msg("starting container")
	startedAt := r.clock.Now()
	if err := r.dockerClient.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	if err := stream.Send(&runnerv1.SubmitRunRequest{
		Id:        r.id,
		RunId:     run.Id,
		StartedAt: timestamppb.New(startedAt),
	}); err != nil {
		return fmt.Errorf("submitting run data : %w", err)
	}
	runLogger.Debug().Msg("submitted started at")

	logErrCh := make(chan error)
	go func() {
		defer close(logErrCh)

		logs, err := r.dockerClient.ContainerLogs(ctx, createRes.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Follow:     true,
			Tail:       "all",
		})
		if err != nil {
			logErrCh <- err
			return
		}
		defer logs.Close()

		stdOutStreamer := newRunLogStreamPipe(commonv1.Run_Log_STDOUT, stream)
		stdErrStreamer := newRunLogStreamPipe(commonv1.Run_Log_STDERR, stream)

		stdOutStreamer.interceptor = func(srr *runnerv1.SubmitRunRequest, line []byte) {
			m := map[string]string{}
			i := 0

			for i < len(line) {
				start := bytes.Index(line[i:], recordStart)
				if start < 0 {
					break
				}
				i += start + len(recordStart)

				terminator := bytes.Index(line[i:], recordEnd)
				if terminator < 0 {
					break
				}

				record := line[i : terminator+i]

				splitted := strings.SplitN(string(record), ":", 2)

				if len(splitted) == 2 {
					m[splitted[0]] = splitted[1]
				}

				i += terminator + len(recordEnd)
			}

			if len(m) > 0 {
				srr.ResultData = m
			}
		}

		var eg errgroup.Group
		eg.Go(func() error {
			return stdOutStreamer.Stream(stream.Context())
		})
		eg.Go(func() error {
			return stdErrStreamer.Stream(stream.Context())
		})
		eg.Go(func() error {
			// NOTE There's something wrong with stdcopy.StdCopy. Trying to copy the
			// logs directly using it into the stream pipe writer will intermittently
			// hang. Passing it through an io.Pipe seems to fix this.
			// TODO figure out why and get rid of the use of the extra pipe here.
			r, w := io.Pipe()
			go func() {
				defer w.Close()
				if _, err := io.Copy(w, logs); err != nil {
					runLogger.Error().Err(err).Msg("failed to close logs pipe writer")
				}
			}()

			if _, err = stdcopy.StdCopy(stdOutStreamer, stdErrStreamer, r); err != nil {
				return err
			}
			if err := stdOutStreamer.Close(); err != nil {
				runLogger.Error().Err(err).Msg("failed to close stdout stream")
			}
			if err := stdErrStreamer.Close(); err != nil {
				runLogger.Error().Err(err).Msg("failed to close stderr stream")
			}
			return nil
		})
		err = eg.Wait()
		if err != nil {
			logErrCh <- fmt.Errorf("parsing logs: %w", err)
		}
	}()

	statusCh, errCh := r.dockerClient.ContainerWait(ctx, createRes.ID, container.WaitConditionNotRunning)
	select {
	case status := <-statusCh:
		err := <-logErrCh
		now := r.clock.Now()

		runLogger.Info().
			Int64("status_code", status.StatusCode).
			Err(err).
			Msg("run execution completed")

		if err != nil {
			if err := stream.Send(&runnerv1.SubmitRunRequest{
				Result: commonv1.Run_ERROR,
				Logs: []*commonv1.Run_Log{
					{
						Time:       now.Format(time.RFC3339Nano),
						OutputType: commonv1.Run_Log_TSTR,
						Data:       []byte(fmt.Sprintf("test run execution failed to capture test output: %s", err)),
					},
				},
				FinishedAt: timestamppb.New(now),
			}); err != nil {
				runLogger.Error().Err(err).Msg("failed submitting error run due to log collection")
			}
			return fmt.Errorf("streaming run execution logs: %w", err)
		}

		var result commonv1.Run_Result
		if status.StatusCode == 0 {
			result = commonv1.Run_PASS
		} else {
			result = commonv1.Run_FAIL
		}
		if err := stream.Send(&runnerv1.SubmitRunRequest{
			Result:     result,
			FinishedAt: timestamppb.New(now),
		}); err != nil {
			return fmt.Errorf("submitting completed run data : %w", err)
		}
		runLogger.Debug().Msg("submitted finshed at")

		runLogger.Info().Msg("completed run")
		return nil
	case err := <-errCh:
		now := r.clock.Now()
		if err := stream.Send(&runnerv1.SubmitRunRequest{
			Result: commonv1.Run_ERROR,
			Logs: []*commonv1.Run_Log{
				{
					Time:       now.Format(time.RFC3339Nano),
					OutputType: commonv1.Run_Log_TSTR,
					Data:       []byte(fmt.Sprintf("test run execution failed: %s", err)),
				},
			},
			FinishedAt: timestamppb.New(now),
		}); err != nil {
			runLogger.Error().Err(err).Msg("failed submitting error run due to execution error")
		}
		return fmt.Errorf("executing run: %w", err)
	}
}

type runLogStreamPipe struct {
	stream      runnerv1.RunnerService_SubmitRunClient
	outputType  commonv1.Run_Log_Output
	scanner     *bufio.Scanner
	writer      io.WriteCloser
	interceptor func(*runnerv1.SubmitRunRequest, []byte)
}

func newRunLogStreamPipe(outputType commonv1.Run_Log_Output, stream runnerv1.RunnerService_SubmitRunClient) *runLogStreamPipe {
	r, w := io.Pipe()
	return &runLogStreamPipe{
		stream:     stream,
		outputType: outputType,
		scanner:    bufio.NewScanner(r),
		writer:     w,
	}
}

func (w *runLogStreamPipe) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *runLogStreamPipe) Close() error {
	return w.writer.Close()
}

func (w *runLogStreamPipe) Stream(ctx context.Context) error {
	for w.scanner.Scan() {
		line := w.scanner.Bytes()

		time, data, found := bytes.Cut(line, []byte(" "))
		if !found {
			// Ignore time if not found
			time, data = []byte(""), time
		}

		submitRunRequest := &runnerv1.SubmitRunRequest{
			Logs: []*commonv1.Run_Log{{
				Time:       string(time),
				OutputType: w.outputType,
				Data:       data,
			}},
		}

		if w.interceptor != nil {
			w.interceptor(submitRunRequest, line)
		}

		if err := w.stream.Send(submitRunRequest); err != nil {
			return fmt.Errorf("streaming log: %w", err)
		}
	}
	if err := w.scanner.Err(); err != nil {
		return fmt.Errorf("reading from %s: %w", w.outputType.String(), err)
	}
	return nil
}
