import axios from 'axios';
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(duration)
dayjs.extend(relativeTime)

import  './api/data/v1/data.pb';
import { DataService, GetRunnerRequest } from './api/data/v1/data.pb';
import { InitReq } from './api/fetch.pb';

const initReq : InitReq = {
  pathPrefix: '/api'
}

const APINullableToJs = function (obj: any, root: boolean = true) : any {
  if (obj == null) {
    return null;
  }

  const datetimeFields = new Set()
  datetimeFields.add('finishedAt')
  datetimeFields.add('nextRunAt')
  datetimeFields.add('registeredAt')
  datetimeFields.add('scheduledAt')
  datetimeFields.add('startedAt')
  datetimeFields.add('updatedAt')
  datetimeFields.add('lastHeartbeatAt')

  if (root) {
    if (Array.isArray(obj)) {
      return obj.map((o) => APINullableToJs(o, true));
    } else {
      Object.keys(obj).map(function (key) {
        if (datetimeFields.has(key)) {
          obj[key] = obj[key] != null? dayjs(obj[key]) : null
        }
      });
      return obj;
    }
  }

  const len = Object.keys(obj).length;

  if (len == 2 && "Time" in obj && "Valid" in obj) {
    if (!obj.Valid) return null;
    return dayjs(obj.Time);
  }

  if (len == 2 && "String" in obj && "Valid" in obj) {
    if (!obj.Valid) return null;
    return obj.String;
  }

  return obj;
};

export default {
  fetchTests: async function (
    ids: string[] = []
  ) {
    const url = "/api/data/v1/tests";
    const params = new URLSearchParams();
    ids.forEach(id => params.append("ids", id));

    return await axios.get(`${url}?${params.toString()}`).then(r => r.data.tests)
  },

  fetchRunners: async function (
    ids: string[] = []
  ) {
    const url = "/api/data/v1/runners";
    const params = new URLSearchParams();
    ids.forEach(id => params.append("ids", id));

    return await axios.get(`${url}?${params.toString()}`).then(r => APINullableToJs(r.data.runners))
  },

  fetchRunDetails: async function (runId: String) {
    const url = `/api/data/v1/runs/${runId}`;
    return await axios.get(url).then (r => APINullableToJs(r.data.run))
  },

  fetchRuns: async function () {
    const url = `/api/data/v1/runs`;
    return await axios.get(url).then (r => APINullableToJs(r.data.runs))
  },

  fetchTestDetails: async function (
    testId: string,
    includeRuns: boolean = true
  ) {
    const url = `/api/data/v1/tests/${testId}?runs=${includeRuns ? 100 : 0}`;

    var testDetails = await axios.get(url).then(r => r.data)
    // console.log("API_",testDetails)

    testDetails.test = APINullableToJs(testDetails.test)
    testDetails.runSummaries = APINullableToJs(testDetails.runSummaries)

    // testDetails = APINullableToJs(testDetails);
    // testDetails.RunsSummary = APINullableToJs(testDetails.RunsSummary);

    return testDetails;
  },

  fetchRunnerDetails: async function (
    runnerId: string,
    includeRuns: boolean = true
  ) {
    const url = `/api/data/v1/runners/${runnerId}?runs=${includeRuns ? 100 : 0}`;
    var data = await axios.get(url).then(r => r.data)

    data.runner = APINullableToJs(data.runner);

    let z = await DataService.GetRunner({id: runnerId},initReq)
    console.log(z)

    if (data.runSummaries != null) {
      data.runSummaries = data.runSummaries.map(function (run:any) {
        run = APINullableToJs(run);
        return run;
      });
    }
    return data;
  },
};

