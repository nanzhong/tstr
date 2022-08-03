<script setup lang="ts">
import { inject, ref, onMounted } from "vue";
import dayjs from "dayjs";
import Plotly from "plotly.js-dist/plotly";
import { LightningBoltIcon, EmojiHappyIcon, EmojiSadIcon, ExclamationCircleIcon } from '@heroicons/vue/outline';
import { RunResult } from "../api/common/v1/common.pb";
import { DataService, SummarizeRunsRequestInterval } from "../api/data/v1/data.pb";
import { InitReq } from "../api/fetch.pb";
import RunsTable from "../components/RunsTable.vue";
import TestRunCountsTable from "../components/TestRunCountsTable.vue";

const initReq: InitReq = inject('dataInitReq')!;

const now = dayjs();
const dayAgoCeil = now.startOf("hour").subtract(23, "hour");
const runs = (await DataService.QueryRuns({ scheduledAfter: dayAgoCeil.toISOString() }, initReq)).runs;
const runSummary = (await DataService.SummarizeRuns({
  scheduledAfter: now.subtract(1, "day").toISOString(),
  window: `${24*60*60}s`,
  interval: SummarizeRunsRequestInterval.HOUR,
}, initReq));

const hourBuckets: string[] = [];
const runDataSeries = {
  [RunResult.UNKNOWN]: { name: "Pending", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#0ea5e9" } },
  [RunResult.FAIL]: { name: "Fail", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ef4444" } },
  [RunResult.ERROR]: { name: "Error", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ec4899" } },
  [RunResult.PASS]: { name: "Pass", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#22c55e" } },
}

const runCounts = {
  [RunResult.UNKNOWN]: 0,
  [RunResult.FAIL]: 0,
  [RunResult.ERROR]: 0,
  [RunResult.PASS]: 0,
}

const testRunCounts: { [key: string]: { testID: string, testName: string, runCounts: { [key in RunResult]: number } } } = {};

runSummary.intervalStats?.forEach(s => {
  hourBuckets.push(dayjs(s.startTime!).toISOString());
  s.resultCount?.forEach(c => {
    switch(c.result!) {
      case RunResult.PASS:
      case RunResult.FAIL:
      case RunResult.ERROR:
      case RunResult.UNKNOWN:
        runDataSeries[c.result!].y.push(c.count || 0);
        runCounts[c.result!] += (c.count! || 0);
      break;
      default:
        runDataSeries[RunResult.UNKNOWN].y.push(c.count || 0);
        runCounts[RunResult.UNKNOWN] += (c.count! || 0);
    }
  });

  s.testCount?.forEach(c => {
    if (!testRunCounts[c.testId!]) {
      testRunCounts[c.testId!] = {
        testID: c.testId!,
        testName: c.testName!,
        runCounts: {
          [RunResult.UNKNOWN]: 0,
          [RunResult.FAIL]: 0,
          [RunResult.ERROR]: 0,
          [RunResult.PASS]: 0,
        }
      }
    }

    c.resultCount?.forEach(rc => {
      switch(rc.result!) {
        case RunResult.PASS:
        case RunResult.FAIL:
        case RunResult.ERROR:
        case RunResult.UNKNOWN:
          testRunCounts[c.testId!].runCounts[rc.result!] += (rc.count! || 0);
          break;
        default:
        testRunCounts[c.testId!].runCounts[RunResult.UNKNOWN] += (rc.count! || 0);
      }
    });
  });
});

const totalRunCount = runCounts[RunResult.PASS] + runCounts[RunResult.FAIL] + runCounts[RunResult.ERROR] + runCounts[RunResult.UNKNOWN];
const stats = {
  totalRuns: {
    name: 'Runs Scheduled',
    icon: LightningBoltIcon,
    stat: totalRunCount,
    colour: 'text-indigo-500'
  },
  passRate: {
    name: 'Pass Rate',
    icon: EmojiHappyIcon,
    stat: `${(runCounts[RunResult.PASS] / totalRunCount * 100).toFixed(2)}%`,
    colour: 'text-green-500'
  },
  FailRate: {
    name: 'Fail Rate',
    icon: EmojiSadIcon,
    stat: `${(runCounts[RunResult.FAIL] / totalRunCount * 100).toFixed(2)}%`,
    colour: 'text-red-500'
  },
  errorRate: {
    name: 'Error Rate',
    icon: ExclamationCircleIcon,
    stat: `${(runCounts[RunResult.ERROR] / totalRunCount * 100).toFixed(2)}%`,
    colour: 'text-pink-500'
  },
}

var data = [
  runDataSeries[RunResult.PASS],
  runDataSeries[RunResult.FAIL],
  runDataSeries[RunResult.ERROR],
  runDataSeries[RunResult.UNKNOWN],
];

var layout = {
  barmode: 'stack',
  bargap: 0.05,
  height: 200,
  margin: {
    l: 40,
    r: 40,
    t: 10,
    b: 40,
  },
};

const recentRunsPlot = ref(null);

onMounted(() => {
  Plotly.newPlot(recentRunsPlot.value, data, layout, { repsonsive: true, displayModeBar: false });
})
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Summary</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <dl class="mt-5 grid grid-cols-1 rounded-lg bg-white overflow-hidden shadow divide-y divide-gray-200 md:grid-cols-4 md:divide-y-0 md:divide-x">
        <div v-for="item in stats" :key="item.name" class="px-4 py-5 sm:p-6">
          <dt>
            <div class="absolute bg-indigo-500 rounded-md p-3">
              <component :is="item.icon" class="h-6 w-6 text-white" aria-hidden="true" />
            </div>
            <p class="ml-16 text-sm font-medium text-gray-500">{{ item.name }}</p>
          </dt>
          <dd class="ml-16 flex items-baseline">
            <p :class="[ item.colour ? item.colour: '', 'text-2xl font-semibold']">
              {{ item.stat }}
            </p>
          </dd>
        </div>
      </dl>
    </div>

    <div class="mt-5 bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Runs</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <div ref="recentRunsPlot"></div>

      <div class="mt-2">
        <h3 class="text-md font-medium text-gray-900">Last 5 Runs</h3>
        <div class="mt-5 flex flex-col">
          <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-6">
            <div class="inline-block min-w-full py-2 align-middle">
              <div class="overflow-hidden shadow-sm ring-1 ring-black ring-opacity-5">
                <RunsTable :runs="runs!.slice(0, 5)" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-5 bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Tests</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <div class="mt-5 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-6">
          <div class="inline-block min-w-full py-2 align-middle">
            <div class="overflow-hidden shadow-sm ring-1 ring-black ring-opacity-5">
              <TestRunCountsTable :testRunCounts="Object.values(testRunCounts)" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
