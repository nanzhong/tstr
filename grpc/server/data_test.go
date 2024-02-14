package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	datav1 "github.com/nanzhong/tstr/api/data/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func getContext(namespace string) context.Context {
	if namespace == "" {
		return context.Background()
	} else {
		md := metadata.New(map[string]string{"namespace": namespace})
		return metadata.NewIncomingContext(context.Background(), md)
	}
}

func createTestData(id string, namespace string) *commonv1.Test {
	testName := "test-name-1"
	timestamp := types.ToProtoTimestamp(time.Now())
	return &commonv1.Test{
		Id:           id,
		Name:         testName,
		Namespace:    namespace,
		RunConfig:    getTestRunConfig(),
		Labels:       map[string]string{"region": "nyc"},
		CronSchedule: "*/15 * * * *",
		NextRunAt:    timestamp,
		Matrix:       &commonv1.Test_Matrix{Labels: map[string]*commonv1.Test_Matrix_LabelValues{}},
		RegisteredAt: timestamp,
		UpdatedAt:    timestamp,
	}
}

func getTestRunConfig() *commonv1.Test_RunConfig {
	return &commonv1.Test_RunConfig{
		ContainerImage: "container-image",
		Command:        "docker",
		Args:           []string{"arg1", "arg2"},
		Env:            map[string]string{"region": "nyc"},
		Timeout:        new(durationpb.Duration),
	}
}

func newTestDataServer(t *testing.T) (*DataServer, *db.MockQuerier) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockQuerier := db.NewMockQuerier(ctrl)

	return &DataServer{
		dbQuerier: mockQuerier,
		clock:     clock.NewMock(),
	}, mockQuerier
}

func TestDataServer_GetTest(t *testing.T) {
	tests := []struct {
		name                           string
		responseCode                   codes.Code
		errMsg                         string
		test                           *commonv1.Test
		mockQuerierReturnTest          func(*commonv1.Test) (db.Test, error)
		mockQuerierReturnTestSummaries func(*commonv1.Test) ([]db.RunSummariesForTestRow, error)
	}{
		{
			name:         "get test using valid test data",
			responseCode: codes.OK,
			errMsg:       "",
			test:         createTestData(uuid.New().String(), "ns1"),
			mockQuerierReturnTest: func(test *commonv1.Test) (db.Test, error) {
				testId, _ := uuid.Parse(test.Id)
				dbTest := db.Test{
					ID:           testId,
					Name:         test.Name,
					CronSchedule: sql.NullString{Valid: true, String: test.CronSchedule},
					NextRunAt:    types.FromProtoTimestampAsNullTime(test.NextRunAt),
					RegisteredAt: types.FromProtoTimestampAsNullTime(test.RegisteredAt),
					UpdatedAt:    types.FromProtoTimestampAsNullTime(test.UpdatedAt),
					Namespace:    test.Namespace,
				}
				dbTest.RunConfig.Set(test.RunConfig)
				dbTest.Labels.Set(test.Labels)
				dbTest.Matrix.Set(test.Matrix)
				return dbTest, nil
			},
			mockQuerierReturnTestSummaries: func(test *commonv1.Test) ([]db.RunSummariesForTestRow, error) {
				return []db.RunSummariesForTestRow{}, nil
			},
		},
		{
			name:         "get test with no namespace in context",
			responseCode: codes.InvalidArgument,
			errMsg:       "request metadata missing namespace",
			test:         createTestData(uuid.New().String(), ""),
			mockQuerierReturnTest: func(test *commonv1.Test) (db.Test, error) {
				return db.Test{}, nil
			},
			mockQuerierReturnTestSummaries: func(test *commonv1.Test) ([]db.RunSummariesForTestRow, error) {
				return []db.RunSummariesForTestRow{}, nil
			},
		},
		{
			name:         "get test using invalid test id",
			responseCode: codes.InvalidArgument,
			errMsg:       "invalid test id",
			test:         createTestData("invalid", "ns1"),
			mockQuerierReturnTest: func(test *commonv1.Test) (db.Test, error) {
				return db.Test{}, nil
			},
			mockQuerierReturnTestSummaries: func(test *commonv1.Test) ([]db.RunSummariesForTestRow, error) {
				return []db.RunSummariesForTestRow{}, nil
			},
		},
		{
			name:         "fail query for get test",
			responseCode: codes.Internal,
			errMsg:       "failed to get test",
			test:         createTestData(uuid.New().String(), "ns1"),
			mockQuerierReturnTest: func(test *commonv1.Test) (db.Test, error) {
				return db.Test{}, errors.New("Dummy Error")
			},
			mockQuerierReturnTestSummaries: func(test *commonv1.Test) ([]db.RunSummariesForTestRow, error) {
				return []db.RunSummariesForTestRow{}, nil
			},
		},
	}

	server, mockQuerier := newTestDataServer(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getContext(tt.test.Namespace)
			testId, _ := uuid.Parse(tt.test.Id)

			if tt.responseCode != codes.InvalidArgument {
				mockQuerier.EXPECT().GetTest(ctx, gomock.AssignableToTypeOf(server.pgxPool), db.GetTestParams{
					ID:        testId,
					Namespace: tt.test.Namespace,
				}).Return(tt.mockQuerierReturnTest(tt.test))

				if _, err := tt.mockQuerierReturnTest(tt.test); err == nil {
					mockQuerier.EXPECT().RunSummariesForTest(ctx, gomock.AssignableToTypeOf(server.pgxPool), db.RunSummariesForTestParams{
						TestID:         testId,
						Namespace:      tt.test.Namespace,
						ScheduledAfter: sql.NullTime{Valid: true, Time: server.clock.Now().Add(-24 * time.Hour)},
					}).Return(tt.mockQuerierReturnTestSummaries(tt.test))
				}
			}

			testRequest := datav1.GetTestRequest{
				Id: tt.test.Id,
			}
			res, err := server.GetTest(ctx, &testRequest)

			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &datav1.GetTestResponse{Test: tt.test}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, tt.responseCode, er.Code())
					assert.Equal(t, tt.errMsg, er.Message())
				}
			}
		})
	}
}

func TestDataServer_QueryTests(t *testing.T) {
	tests := []struct {
		name                   string
		responseCode           codes.Code
		errMsg                 string
		namespace              string
		tests                  func(int, string) []*commonv1.Test
		mockQuerierReturnTests func([]*commonv1.Test) ([]db.Test, error)
	}{
		{
			name:         "query tests using valid test data",
			responseCode: codes.OK,
			errMsg:       "",
			namespace:    "ns1",
			tests: func(testCount int, namespace string) []*commonv1.Test {
				var tests []*commonv1.Test
				for i := 1; i <= testCount; i++ {
					test := createTestData(uuid.New().String(), namespace)
					tests = append(tests, test)
				}
				return tests
			},
			mockQuerierReturnTests: func(tests []*commonv1.Test) ([]db.Test, error) {
				var testsToReturn []db.Test
				for _, test := range tests {
					testId, _ := uuid.Parse(test.Id)
					dbTest := db.Test{
						ID:           testId,
						Name:         test.Name,
						CronSchedule: sql.NullString{Valid: true, String: test.CronSchedule},
						NextRunAt:    types.FromProtoTimestampAsNullTime(test.NextRunAt),
						RegisteredAt: types.FromProtoTimestampAsNullTime(test.RegisteredAt),
						UpdatedAt:    types.FromProtoTimestampAsNullTime(test.UpdatedAt),
						Namespace:    test.Namespace,
					}
					dbTest.RunConfig.Set(test.RunConfig)
					dbTest.Labels.Set(test.Labels)
					dbTest.Matrix.Set(test.Matrix)
					testsToReturn = append(testsToReturn, dbTest)
				}
				return testsToReturn, nil
			},
		},
		{
			name:         "query tests with no namespace in context",
			responseCode: codes.InvalidArgument,
			errMsg:       "request metadata missing namespace",
			namespace:    "",
			tests: func(testCount int, namespace string) []*commonv1.Test {
				return []*commonv1.Test{}
			},
			mockQuerierReturnTests: func(tests []*commonv1.Test) ([]db.Test, error) {
				return []db.Test{}, nil
			},
		},
		{
			name:         "query tests with invalid test id",
			responseCode: codes.InvalidArgument,
			errMsg:       "failed to parse test id",
			namespace:    "ns1",
			tests: func(testCount int, namespace string) []*commonv1.Test {
				var tests []*commonv1.Test
				for i := 1; i <= testCount; i++ {
					test := createTestData("invalid", namespace)
					tests = append(tests, test)
				}
				return tests
			},
			mockQuerierReturnTests: func(tests []*commonv1.Test) ([]db.Test, error) {
				return []db.Test{}, nil
			},
		},
	}

	server, mockQuerier := newTestDataServer(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getContext(tt.namespace)
			tests := tt.tests(2, tt.namespace)

			var testUUIDs []uuid.UUID
			var testIds []string
			for _, test := range tests {
				testId, _ := uuid.Parse(test.Id)
				testUUIDs = append(testUUIDs, testId)
				testIds = append(testIds, test.Id)
			}

			var labels map[string]string
			var byteLabels []byte

			if len(tests) > 0 {
				rand.Seed(time.Now().Unix())
				labels = tests[rand.Intn(len(tests))].Labels
				byteLabels, _ = json.Marshal(labels)
			}

			if tt.responseCode != codes.InvalidArgument {
				mockQuerier.EXPECT().QueryTests(ctx, gomock.AssignableToTypeOf(server.pgxPool), db.QueryTestsParams{
					Namespace: tt.namespace,
					Ids:       testUUIDs,
					Labels:    pgtype.JSONB{Bytes: byteLabels, Status: pgtype.Present},
				}).Return(tt.mockQuerierReturnTests(tests))
			}

			queryTestsRequest := datav1.QueryTestsRequest{
				Ids:          testIds,
				TestSuiteIds: []string{},
				Labels:       labels,
			}
			res, err := server.QueryTests(ctx, &queryTestsRequest)

			if res != nil {
				require.NoError(t, err)
				assert.Equal(t, &datav1.QueryTestsResponse{Tests: tests}, res)
			} else {
				if er, ok := status.FromError(err); ok {
					assert.Equal(t, tt.responseCode, er.Code())
					assert.Equal(t, tt.errMsg, er.Message())
				}
			}
		})
	}
}
