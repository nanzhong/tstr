<script setup lang="ts">
import { inject, ref, onMounted } from "vue";
import dayjs from "dayjs";
import Plotly from "plotly.js-dist/plotly";
import { Run, RunResult } from "../api/common/v1/common.pb";
import { DataService, SummarizeRunsRequestInterval } from "../api/data/v1/data.pb";
import { InitReq } from "../api/fetch.pb";
import RunsTable from "../components/RunsTable.vue";
import { ArrowSmDownIcon, ArrowSmUpIcon } from "@heroicons/vue/solid";

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

runSummary.intervalStats?.forEach(s => {
  hourBuckets.push(dayjs(s.startTime!).toISOString());
  s.resultCount?.forEach(c => {
    switch(c.result!) {
      case RunResult.PASS:
      case RunResult.FAIL:
      case RunResult.ERROR:
      case RunResult.UNKNOWN:
        runDataSeries[c.result!].y.push(c.count || 0);
      break;
      default:
        runDataSeries[RunResult.UNKNOWN].y.push(c.count || 0);
    }
  });
});
console.log(runDataSeries);

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

const stats = [
  { name: 'Test Runs Orchestrated', stat: '71,897', colour: 'text-indigo-500' },
  { name: 'Pass Rate', stat: '58.16%', colour: 'text-green-500' },
  { name: 'Error Rate', stat: '24.57%', colour: 'text-pink-500' },
]

onMounted(() => {
  Plotly.newPlot(recentRunsPlot.value, data, layout, { repsonsive: true, displayModeBar: false });
})
</script>

<template>
  <div>
    <div class="pb-3 border-b border-gray-500">
      <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
        <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Summary</h2>
        <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
      </div>
    </div>

    <dl
      class="mt-5 grid grid-cols-1 rounded-lg bg-white overflow-hidden shadow divide-y divide-gray-200 md:grid-cols-3 md:divide-y-0 md:divide-x">
      <div v-for="item in stats" :key="item.name" class="px-4 py-5 sm:p-6">
        <dt class="text-base font-normal text-gray-900">
          {{ item.name }}
        </dt>
        <dd class="mt-1 flex justify-between items-baseline md:block lg:flex">
          <div :class="[ item.colour ? item.colour: '', 'flex items-baseline text-2xl font-semibold']">
            {{ item.stat }}
          </div>
        </dd>
      </div>
    </dl>
  </div>

  <div>
    <div class="mt-5 pb-3 border-b border-gray-500">
      <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
        <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Recent Test Runs</h2>
        <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
      </div>
    </div>

    <div ref="recentRunsPlot"></div>

    <div class="mt-2">
      <h3 class="text-lg font-medium text-gray-900">Last 5 Test Runs</h3>
      <div class="mt-5 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
            <RunsTable :runs="runs!.slice(0, 5)" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
