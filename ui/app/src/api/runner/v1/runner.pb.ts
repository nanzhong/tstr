/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"
export type RegisterRunnerRequest = {
  name?: string
  acceptTestLabelSelectors?: {[key: string]: string}
  rejectTestLabelSelectors?: {[key: string]: string}
}

export type RegisterRunnerResponse = {
  runner?: TstrCommonV1Common.Runner
}

export type NextRunRequest = {
  id?: string
}

export type NextRunResponse = {
  run?: TstrCommonV1Common.Run
}

export type SubmitRunRequest = {
  id?: string
  runId?: string
  result?: TstrCommonV1Common.RunResult
  logs?: TstrCommonV1Common.RunLog[]
  resultData?: {[key: string]: string}
  startedAt?: GoogleProtobufTimestamp.Timestamp
  finishedAt?: GoogleProtobufTimestamp.Timestamp
}

export type SubmitRunResponse = {
}

export class RunnerService {
  static RegisterRunner(req: RegisterRunnerRequest, initReq?: fm.InitReq): Promise<RegisterRunnerResponse> {
    return fm.fetchReq<RegisterRunnerRequest, RegisterRunnerResponse>(`/tstr.runner.v1.RunnerService/RegisterRunner`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static NextRun(req: NextRunRequest, initReq?: fm.InitReq): Promise<NextRunResponse> {
    return fm.fetchReq<NextRunRequest, NextRunResponse>(`/tstr.runner.v1.RunnerService/NextRun`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}