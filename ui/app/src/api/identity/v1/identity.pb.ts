/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
export type IdentityRequest = {
}

export type IdentityResponse = {
  accessToken?: TstrCommonV1Common.AccessToken
  accessibleNamespaces?: string[]
}

export class IdentityService {
  static Identity(req: IdentityRequest, initReq?: fm.InitReq): Promise<IdentityResponse> {
    return fm.fetchReq<IdentityRequest, IdentityResponse>(`/identity/v1/identity?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}