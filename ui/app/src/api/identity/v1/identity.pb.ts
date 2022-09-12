/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"
export type IdentityRequest = {
}

export type IdentityResponse = {
  scopes?: TstrCommonV1Common.AccessTokenScope[]
  namespaceSelectors?: string[]
  issuedAt?: GoogleProtobufTimestamp.Timestamp
  expiresAt?: GoogleProtobufTimestamp.Timestamp
}

export class IdentityService {
  static Identity(req: IdentityRequest, initReq?: fm.InitReq): Promise<IdentityResponse> {
    return fm.fetchReq<IdentityRequest, IdentityResponse>(`/identity/v1/identity?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}