/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as GoogleProtobufDuration from "../../google/protobuf/duration.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum RunResult {
  UNKNOWN = "UNKNOWN",
  PASS = "PASS",
  FAIL = "FAIL",
  ERROR = "ERROR",
}

export enum RunLogOutput {
  UNKNOWN = "UNKNOWN",
  STDOUT = "STDOUT",
  STDERR = "STDERR",
  TSTR = "TSTR",
}

export enum AccessTokenScope {
  UNKNOWN = "UNKNOWN",
  ADMIN = "ADMIN",
  CONTROL_R = "CONTROL_R",
  CONTROL_RW = "CONTROL_RW",
  RUNNER = "RUNNER",
  DATA = "DATA",
}

export type TestRunConfig = {
  containerImage?: string
  command?: string
  args?: string[]
  env?: {[key: string]: string}
  timeout?: GoogleProtobufDuration.Duration
}

export type TestMatrixLabelValues = {
  values?: string[]
}

export type TestMatrix = {
  labels?: {[key: string]: TestMatrixLabelValues}
}

export type Test = {
  id?: string
  namespace?: string
  name?: string
  runConfig?: TestRunConfig
  labels?: {[key: string]: string}
  cronSchedule?: string
  nextRunAt?: GoogleProtobufTimestamp.Timestamp
  matrix?: TestMatrix
  registeredAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}

export type RunLog = {
  time?: string
  outputType?: RunLogOutput
  data?: Uint8Array
}

export type Run = {
  id?: string
  testId?: string
  testRunConfig?: TestRunConfig
  testMatrixId?: string
  labels?: {[key: string]: string}
  runnerId?: string
  result?: RunResult
  logs?: RunLog[]
  resultData?: {[key: string]: string}
  scheduledAt?: GoogleProtobufTimestamp.Timestamp
  startedAt?: GoogleProtobufTimestamp.Timestamp
  finishedAt?: GoogleProtobufTimestamp.Timestamp
}

export type Runner = {
  id?: string
  name?: string
  namespaceSelectors?: string[]
  acceptTestLabelSelectors?: {[key: string]: string}
  rejectTestLabelSelectors?: {[key: string]: string}
  registeredAt?: GoogleProtobufTimestamp.Timestamp
  lastHeartbeatAt?: GoogleProtobufTimestamp.Timestamp
}

export type AccessToken = {
  id?: string
  name?: string
  token?: string
  namespaceSelectors?: string[]
  scopes?: AccessTokenScope[]
  issuedAt?: GoogleProtobufTimestamp.Timestamp
  expiresAt?: GoogleProtobufTimestamp.Timestamp
  revokedAt?: GoogleProtobufTimestamp.Timestamp
}