import { DateTime } from "luxon";
import axios from 'axios';

const APINullableToJs = function (obj: any, root: boolean = true) {
  if (obj == null) {
    return null;
  }

  if (root) {
    if (Array.isArray(obj)) {
      return obj.map((o) => APINullableToJs(o, true));
    } else {
      Object.keys(obj).map(function (key) {
        obj[key] = APINullableToJs(obj[key], false);
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
  fetchTests: async function () {
    const url = "/api/tests";
    return await axios.get(url).then(r => r.data)
  },

  fetchRunners: async function () {
    const url = "/api/runners";
    return await axios.get(url).then(r => r.data)
  },

  fetchRunDetails: async function (runId: String) {
    const url = `/api/runs/${runId}`;
    return await axios.get(url).then (r => APINullableToJs(r.data))
  },

  fetchRuns: async function () {
    const url = `/api/runs`;
    return await axios.get(url).then (r => APINullableToJs(r.data))
  },

  fetchTestDetails: async function (
    testId: String,
    includeRuns: boolean = true
  ) {
    const url = `/api/tests/${testId}?runs=${includeRuns ? 100 : 0}`;

    var testDetails = await axios.get(url).then(r => r.data)

    testDetails = APINullableToJs(testDetails);
    testDetails.RunsSummary = APINullableToJs(testDetails.RunsSummary);

    return testDetails;
  },

  fetchRunnerDetails: async function (
    runnerId: String,
    includeRuns: boolean = true
  ) {
    const url = `/api/runners/${runnerId}?runs=${includeRuns ? 100 : 0}`;
    var data = await axios.get(url).then(r => r.data)

    var runner = data.Runner;

    runner = APINullableToJs(data.Runner);

    if (data.RunsSummary != null) {
      runner["LastRuns"] = data.RunsSummary.map(function (run) {
        run = APINullableToJs(run);

        return run;
      });
    }
    return runner;
  },
};
