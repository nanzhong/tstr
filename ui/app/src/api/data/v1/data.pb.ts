/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufDuration from "../../google/protobuf/duration.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum SummarizeRunsRequestInterval {
  UNKNOWN = "UNKNOWN",
  HOUR = "HOUR",
  DAY = "DAY",
  WEEK = "WEEK",
}

export type GetTestRequest = {
  id?: string
}

export type GetTestResponse = {
  test?: TstrCommonV1Common.Test
  runSummaries?: RunSummary[]
}

export type QueryTestsRequest = {
  ids?: string[]
  testSuiteIds?: string[]
  labels?: {[key: string]: string}
}

export type QueryTestsResponse = {
  tests?: TstrCommonV1Common.Test[]
}

export type GetTestSuiteRequest = {
  id?: string
}

export type GetTestSuiteResponse = {
  testSuite?: TstrCommonV1Common.TestSuite
}

export type QueryTestSuitesRequest = {
  ids?: string[]
  labels?: {[key: string]: string}
}

export type QueryTestSuitesResponse = {
  testSuites?: TstrCommonV1Common.TestSuite[]
}

export type GetRunRequest = {
  id?: string
}

export type GetRunResponse = {
  run?: TstrCommonV1Common.Run
}

export type QueryRunsRequest = {
  ids?: string[]
  testIds?: string[]
  testSuiteIds?: string[]
  runnerIds?: string[]
  results?: TstrCommonV1Common.RunResult[]
  scheduledBefore?: GoogleProtobufTimestamp.Timestamp
  scheduledAfter?: GoogleProtobufTimestamp.Timestamp
  startedBefore?: GoogleProtobufTimestamp.Timestamp
  startedAfter?: GoogleProtobufTimestamp.Timestamp
  finishedBefore?: GoogleProtobufTimestamp.Timestamp
  finishedAfter?: GoogleProtobufTimestamp.Timestamp
  includeLogs?: boolean
}

export type QueryRunsResponse = {
  runs?: TstrCommonV1Common.Run[]
}

export type SummarizeRunsRequest = {
  scheduledAfter?: GoogleProtobufTimestamp.Timestamp
  window?: GoogleProtobufDuration.Duration
  interval?: SummarizeRunsRequestInterval
}

export type SummarizeRunsResponseIntervalStatsResultBreakdown = {
  result?: TstrCommonV1Common.RunResult
  count?: number
}

export type SummarizeRunsResponseIntervalStatsTestBreakdown = {
  testId?: string
  testName?: string
  resultCount?: SummarizeRunsResponseIntervalStatsResultBreakdown[]
}

export type SummarizeRunsResponseIntervalStats = {
  startTime?: GoogleProtobufTimestamp.Timestamp
  duration?: GoogleProtobufDuration.Duration
  resultCount?: SummarizeRunsResponseIntervalStatsResultBreakdown[]
  testCount?: SummarizeRunsResponseIntervalStatsTestBreakdown[]
}

export type SummarizeRunsResponse = {
  intervalStats?: SummarizeRunsResponseIntervalStats[]
}

export type GetRunnerRequest = {
  id?: string
}

export type GetRunnerResponse = {
  runner?: TstrCommonV1Common.Runner
  runSummaries?: RunSummary[]
}

export type QueryRunnersRequest = {
  ids?: string[]
  lastHeartbeatWithin?: GoogleProtobufDuration.Duration
}

export type QueryRunnersResponse = {
  runners?: TstrCommonV1Common.Runner[]
}

export type RunSummary = {
  id?: string
  testId?: string
  testName?: string
  testRunConfig?: TstrCommonV1Common.TestRunConfig
  labels?: {[key: string]: string}
  runnerId?: string
  result?: TstrCommonV1Common.RunResult
  resultData?: {[key: string]: string}
  scheduledAt?: GoogleProtobufTimestamp.Timestamp
  startedAt?: GoogleProtobufTimestamp.Timestamp
  finishedAt?: GoogleProtobufTimestamp.Timestamp
}

export class DataService {
  static GetTest(req: GetTestRequest, initReq?: fm.InitReq): Promise<GetTestResponse> {
    return fm.fetchReq<GetTestRequest, GetTestResponse>(`/data/v1/tests/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static QueryTests(req: QueryTestsRequest, initReq?: fm.InitReq): Promise<QueryTestsResponse> {
    return fm.fetchReq<QueryTestsRequest, QueryTestsResponse>(`/data/v1/tests?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GetTestSuite(req: GetTestSuiteRequest, initReq?: fm.InitReq): Promise<GetTestSuiteResponse> {
    return fm.fetchReq<GetTestSuiteRequest, GetTestSuiteResponse>(`/data/v1/test_suites/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static QueryTestSuites(req: QueryTestSuitesRequest, initReq?: fm.InitReq): Promise<QueryTestSuitesResponse> {
    return fm.fetchReq<QueryTestSuitesRequest, QueryTestSuitesResponse>(`/data/v1/test_suites?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GetRun(req: GetRunRequest, initReq?: fm.InitReq): Promise<GetRunResponse> {
    return fm.fetchReq<GetRunRequest, GetRunResponse>(`/data/v1/runs/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static QueryRuns(req: QueryRunsRequest, initReq?: fm.InitReq): Promise<QueryRunsResponse> {
    return fm.fetchReq<QueryRunsRequest, QueryRunsResponse>(`/data/v1/runs?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static SummarizeRuns(req: SummarizeRunsRequest, initReq?: fm.InitReq): Promise<SummarizeRunsResponse> {
    return fm.fetchReq<SummarizeRunsRequest, SummarizeRunsResponse>(`/data/v1/runs/summary?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GetRunner(req: GetRunnerRequest, initReq?: fm.InitReq): Promise<GetRunnerResponse> {
    return fm.fetchReq<GetRunnerRequest, GetRunnerResponse>(`/data/v1/runners/${req["id"]}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"})
  }
  static QueryRunners(req: QueryRunnersRequest, initReq?: fm.InitReq): Promise<QueryRunnersResponse> {
    return fm.fetchReq<QueryRunnersRequest, QueryRunnersResponse>(`/data/v1/runners?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}