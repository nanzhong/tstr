/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufField_mask from "../../google/protobuf/field_mask.pb"
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

export type DeleteTestRequest = {
  id?: string
}

export type DeleteTestResponse = {
}

export type ScheduleRunRequest = {
  testId?: string
  labels?: {[key: string]: string}
  testMatrix?: TstrCommonV1Common.TestMatrix
}

export type ScheduleRunResponse = {
  runs?: TstrCommonV1Common.Run[]
}

export class ControlService {
  static RegisterTest(req: RegisterTestRequest, initReq?: fm.InitReq): Promise<RegisterTestResponse> {
    return fm.fetchReq<RegisterTestRequest, RegisterTestResponse>(`/tstr.control.v1.ControlService/RegisterTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateTest(req: UpdateTestRequest, initReq?: fm.InitReq): Promise<UpdateTestResponse> {
    return fm.fetchReq<UpdateTestRequest, UpdateTestResponse>(`/tstr.control.v1.ControlService/UpdateTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteTest(req: DeleteTestRequest, initReq?: fm.InitReq): Promise<DeleteTestResponse> {
    return fm.fetchReq<DeleteTestRequest, DeleteTestResponse>(`/tstr.control.v1.ControlService/DeleteTest`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ScheduleRun(req: ScheduleRunRequest, initReq?: fm.InitReq): Promise<ScheduleRunResponse> {
    return fm.fetchReq<ScheduleRunRequest, ScheduleRunResponse>(`/tstr.control.v1.ControlService/ScheduleRun`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}