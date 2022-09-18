<script setup lang="ts">
import dayjs from "dayjs";
import { useInitReq } from "../api/init";
import { DataService, SummarizeRunsRequestInterval } from "../api/data/v1/data.pb";
import { RunResult } from "../api/common/v1/common.pb";
import TestCard from "../components/TestCard.vue";

const initReq = useInitReq();
const now = dayjs();
const runSummary = (await DataService.SummarizeRuns({
  scheduledAfter: now.subtract(1, "day").toISOString(),
  window: `${24*60*60}s`,
  interval: SummarizeRunsRequestInterval.HOUR,
}, initReq));

const hourBuckets = [];
const tests: { [key: string]: { id: string, name: string, results: { interval: dayjs.Dayjs, counts: { [key in RunResult]: number } }[] } } = {};
runSummary.intervalStats?.forEach(s => {
  const interval = dayjs(s.startTime!);
  hourBuckets.push(interval.toISOString());

  s.testCount?.forEach(t => {
    if (!tests[t.testId!]) {
      tests[t.testId!] = {
        id: t.testId!, 
        name: t.testName!,
        results: [],
      };
    }

    const counts = {
      [RunResult.UNKNOWN]: 0,
      [RunResult.FAIL]: 0,
      [RunResult.ERROR]: 0,
      [RunResult.PASS]: 0,
    };
    t.resultCount!.forEach(c => {
      switch(c.result!) {
        case RunResult.PASS:
        case RunResult.FAIL:
        case RunResult.ERROR:
        case RunResult.UNKNOWN:
          counts[c.result!] += (c.count! || 0);
          break;
        default:
          counts[RunResult.UNKNOWN] += (c.count! || 0);
      }
    });

    tests[t.testId!].results.push({
      interval: interval,
      counts: counts,
    });
  });
});

</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Active Tests</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <ul role="list" class="mt-5 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        <li v-for="(test, testID) in tests" :key="testID" class="col-span-1">
          <TestCard :test="test" />
        </li>
      </ul>
    </div>
  </div>
</template>
