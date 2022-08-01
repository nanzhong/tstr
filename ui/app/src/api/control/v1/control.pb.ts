/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufField_mask from "../../google/protobuf/field_mask.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"
export type RegisterTestRequest = {
  name?: string
  labels?: {[key: string]: string}
  runConfig?: TstrCommonV1Common.TestRunConfig
  cronSchedule?: string
  matrix?: TstrCommonV1Common.TestMatrix
}

export type RegisterTestResponse = {
  test?: TstrCommonV1Common.Test
}

export type UpdateTestRequest = {
  fieldMask?: GoogleProtobufField_mask.FieldMask
  id?: string
  name?: string
  runConfig?: TstrCommonV1Common.TestRunConfig
  labels?: {[key: string]: string}
  cronSchedule?: string
  matrix?: TstrCommonV1Common.TestMatrix
}

export type UpdateTestResponse = {
}

export type GetTestRequest = {
  id?: string
}

export type GetTestResponse = {
  test?: TstrCommonV1Common.Test
}

export type ListTestsRequest = {
}

export type ListTestsResponse = {
  tests?: TstrCommonV1Common.Test[]
}

export type ArchiveTestRequest = {
  id?: string
}

export type ArchiveTestResponse = {
}

export type DefineTestSuiteRequest = {
  name?: string
  labels?: {[key: string]: string}
}

export type DefineTestSuiteResponse = {
  testSuite?: TstrCommonV1Common.TestSuite
}

export type UpdateTestSuiteRequest = {
  fieldMask?: GoogleProtobufField_mask.FieldMask
  id?: string
  name?: string
  labels?: {[key: string]: string}
}

export type UpdateTestSuiteResponse = {
}

export type GetTestSuiteRequest = {
  id?: string
}

export type GetTestSuiteResponse = {
  testSuite?: TstrCommonV1Common.TestSuite
}

export type ListTestSuitesRequest = {
}

export type ListTestSuitesResponse = {
  testSuites?: TstrCommonV1Common.TestSuite[]
}

export type ArchiveTestSuiteRequest = {
  id?: string
}

export type ArchiveTestSuiteResponse = {
}

export type GetRunRequest = {
  id?: string
}

export type GetRunResponse = {
  run?: TstrCommonV1Common.Run
}

export type ListRunsRequest = {
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
}

export type ListRunsResponse = {
  runs?: TstrCommonV1Common.Run[]
}

export type ScheduleRunRequest = {
  testId?: string
  labels?: {[key: string]: string}
  testMatrix?: TstrCommonV1Common.TestMatrix
}

export type ScheduleRunResponse = {
  runs?: TstrCommonV1Common.Run[]
}

export type GetRunnerRequest = {
  id?: string
}

export type GetRunnerResponse = {
  runner?: TstrCommonV1Common.Runner
}

export type ListRunnersRequest = {
}

export type ListRunnersResponse = {
  runners?: TstrCommonV1Common.Runner[]
}

export class ControlService {
  static RegisterTest(req: RegisterTestRequest, initReq?: fm.InitReq): Promise<RegisterTestResponse> {
    return fm.fetchReq<RegisterTestRequest, RegisterTestResponse>(`/tstr.control.v1.ControlService/RegisterTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateTest(req: UpdateTestRequest, initReq?: fm.InitReq): Promise<UpdateTestResponse> {
    return fm.fetchReq<UpdateTestRequest, UpdateTestResponse>(`/tstr.control.v1.ControlService/UpdateTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetTest(req: GetTestRequest, initReq?: fm.InitReq): Promise<GetTestResponse> {
    return fm.fetchReq<GetTestRequest, GetTestResponse>(`/tstr.control.v1.ControlService/GetTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListTests(req: ListTestsRequest, initReq?: fm.InitReq): Promise<ListTestsResponse> {
    return fm.fetchReq<ListTestsRequest, ListTestsResponse>(`/tstr.control.v1.ControlService/ListTests`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ArchiveTest(req: ArchiveTestRequest, initReq?: fm.InitReq): Promise<ArchiveTestResponse> {
    return fm.fetchReq<ArchiveTestRequest, ArchiveTestResponse>(`/tstr.control.v1.ControlService/ArchiveTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DefineTestSuite(req: DefineTestSuiteRequest, initReq?: fm.InitReq): Promise<DefineTestSuiteResponse> {
    return fm.fetchReq<DefineTestSuiteRequest, DefineTestSuiteResponse>(`/tstr.control.v1.ControlService/DefineTestSuite`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateTestSuite(req: UpdateTestSuiteRequest, initReq?: fm.InitReq): Promise<UpdateTestSuiteResponse> {
    return fm.fetchReq<UpdateTestSuiteRequest, UpdateTestSuiteResponse>(`/tstr.control.v1.ControlService/UpdateTestSuite`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetTestSuite(req: GetTestSuiteRequest, initReq?: fm.InitReq): Promise<GetTestSuiteResponse> {
    return fm.fetchReq<GetTestSuiteRequest, GetTestSuiteResponse>(`/tstr.control.v1.ControlService/GetTestSuite`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListTestSuites(req: ListTestSuitesRequest, initReq?: fm.InitReq): Promise<ListTestSuitesResponse> {
    return fm.fetchReq<ListTestSuitesRequest, ListTestSuitesResponse>(`/tstr.control.v1.ControlService/ListTestSuites`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ArchiveTestSuite(req: ArchiveTestSuiteRequest, initReq?: fm.InitReq): Promise<ArchiveTestSuiteResponse> {
    return fm.fetchReq<ArchiveTestSuiteRequest, ArchiveTestSuiteResponse>(`/tstr.control.v1.ControlService/ArchiveTestSuite`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetRun(req: GetRunRequest, initReq?: fm.InitReq): Promise<GetRunResponse> {
    return fm.fetchReq<GetRunRequest, GetRunResponse>(`/tstr.control.v1.ControlService/GetRun`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListRuns(req: ListRunsRequest, initReq?: fm.InitReq): Promise<ListRunsResponse> {
    return fm.fetchReq<ListRunsRequest, ListRunsResponse>(`/tstr.control.v1.ControlService/ListRuns`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ScheduleRun(req: ScheduleRunRequest, initReq?: fm.InitReq): Promise<ScheduleRunResponse> {
    return fm.fetchReq<ScheduleRunRequest, ScheduleRunResponse>(`/tstr.control.v1.ControlService/ScheduleRun`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetRunner(req: GetRunnerRequest, initReq?: fm.InitReq): Promise<GetRunnerResponse> {
    return fm.fetchReq<GetRunnerRequest, GetRunnerResponse>(`/tstr.control.v1.ControlService/GetRunner`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListRunners(req: ListRunnersRequest, initReq?: fm.InitReq): Promise<ListRunnersResponse> {
    return fm.fetchReq<ListRunnersRequest, ListRunnersResponse>(`/tstr.control.v1.ControlService/ListRunners`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}