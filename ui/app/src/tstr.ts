import { DateTime } from "luxon";
import axios from 'axios';

const APINullableToJs = function (obj: any, root: boolean = true) {
  if (obj == null) {
    return null;
  }

  const datetime_fields = new Set()
  datetime_fields.add('finished_at')
  datetime_fields.add('next_run_at')
  datetime_fields.add('registered_at')
  datetime_fields.add('scheduled_at')
  datetime_fields.add('started_at')
  datetime_fields.add('updated_at')
  datetime_fields.add('last_heartbeat_at')

  if (root) {
    if (Array.isArray(obj)) {
      return obj.map((o) => APINullableToJs(o, true));
    } else {
      Object.keys(obj).map(function (key) {
        if (datetime_fields.has(key)) {
          obj[key] = obj[key] != null? DateTime.fromISO(obj[key]) : null
        }
      });
      return obj;
    }
  }

  const len = Object.keys(obj).length;

  if (len == 2 && "Time" in obj && "Valid" in obj) {
    if (!obj.Valid) return null;
    return DateTime.fromISO(obj.Time);
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
    testDetails.run_summaries = APINullableToJs(testDetails.run_summaries)

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

    if (data.run_summaries != null) {
      data.run_summaries = data.run_summaries.map(function (run) {
        run = APINullableToJs(run);
        return run;
      });
    }
    return data;
  },
};
