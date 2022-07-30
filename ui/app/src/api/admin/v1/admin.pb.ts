/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as TstrCommonV1Common from "../../common/v1/common.pb"
import * as fm from "../../fetch.pb"
import * as GoogleProtobufDuration from "../../google/protobuf/duration.pb"
export type IssueAccessTokenRequest = {
  name?: string
  scopes?: TstrCommonV1Common.AccessTokenScope[]
  validDuration?: GoogleProtobufDuration.Duration
}

export type IssueAccessTokenResponse = {
  accessToken?: TstrCommonV1Common.AccessToken
}

export type GetAccessTokenRequest = {
  id?: string
}

export type GetAccessTokenResponse = {
  accessToken?: TstrCommonV1Common.AccessToken
}

export type ListAccessTokensRequest = {
  includeExpired?: boolean
  includeRevoked?: boolean
}

export type ListAccessTokensResponse = {
  accessTokens?: TstrCommonV1Common.AccessToken[]
}

export type RevokeAccessTokenRequest = {
  id?: string
}

export type RevokeAccessTokenResponse = {
}

export class AdminService {
  static IssueAccessToken(req: IssueAccessTokenRequest, initReq?: fm.InitReq): Promise<IssueAccessTokenResponse> {
    return fm.fetchReq<IssueAccessTokenRequest, IssueAccessTokenResponse>(`/tstr.admin.v1.AdminService/IssueAccessToken`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetAccessToken(req: GetAccessTokenRequest, initReq?: fm.InitReq): Promise<GetAccessTokenResponse> {
    return fm.fetchReq<GetAccessTokenRequest, GetAccessTokenResponse>(`/tstr.admin.v1.AdminService/GetAccessToken`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListAccessTokens(req: ListAccessTokensRequest, initReq?: fm.InitReq): Promise<ListAccessTokensResponse> {
    return fm.fetchReq<ListAccessTokensRequest, ListAccessTokensResponse>(`/tstr.admin.v1.AdminService/ListAccessTokens`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static RevokeAccessToken(req: RevokeAccessTokenRequest, initReq?: fm.InitReq): Promise<RevokeAccessTokenResponse> {
    return fm.fetchReq<RevokeAccessTokenRequest, RevokeAccessTokenResponse>(`/tstr.admin.v1.AdminService/RevokeAccessToken`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}