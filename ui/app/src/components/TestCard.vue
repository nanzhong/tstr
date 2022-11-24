<script setup lang="ts">
import { ref, onMounted } from "vue";
import dayjs from "dayjs";
import Plotly from "plotly.js-dist/plotly";
import { useInitReq } from "../api/init";
import { DataService, SummarizeRunsRequestInterval } from "../api/data/v1/data.pb";
import { Test } from "../api/common/v1/common.pb";
import { RunResult } from "../api/common/v1/common.pb";

const props = defineProps<{ test: Test }>();

const initReq = useInitReq();
const now = dayjs();
const runSummary = (await DataService.SummarizeRuns({
  scheduledAfter: now.subtract(1, "day").toISOString(),
  window: `${24*60*60}s`,
  interval: SummarizeRunsRequestInterval.HOUR,
  testIds: [props.test.id],
}, initReq));

const plot = ref(null);
const hourBuckets: string[] = [];
const runSeries = {
  [RunResult.PASS]: { name: "Pass", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#22c55e" } },
  [RunResult.FAIL]: { name: "Fail", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ef4444" } },
  [RunResult.ERROR]: { name: "Error", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ec4899" } },
  [RunResult.UNKNOWN]: { name: "Pending", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#0ea5e9" } },
}

if (runSummary.intervalStats) {
  runSummary.intervalStats?.forEach(s => {
    const interval = dayjs(s.startTime!);
    hourBuckets.push(interval.toISOString());

    const counts = {
      [RunResult.UNKNOWN]: 0,
      [RunResult.FAIL]: 0,
      [RunResult.ERROR]: 0,
      [RunResult.PASS]: 0,
    };

    s.resultCount?.forEach(c => {
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

    for (const [result, count] of Object.entries(counts)) {
      runSeries[result].y.push(count);
    }
  });

  const layout = {
    barmode: 'stack',
    bargap: 0.1,
    height: 150,
    margin: {
      l: 20,
      r: 20,
      t: 20,
      b: 40,
    },
    showlegend: false,
  };

  onMounted(() => {
    Plotly.newPlot(plot.value, Object.values(runSeries), layout, { repsonsive: true, displayModeBar: false });
  });
}
</script>

<template>
  <div class="bg-white rounded-lg shadow divide-y divide-gray-200">
    <router-link :to="{ name: 'test-details', params: { id: test.id } }">
      <h3 class="p-5 text-gray-900 text-sm font-medium truncate">{{ test.name }}</h3>
    </router-link>
    <div v-if="runSummary.intervalStats" ref="plot"></div>
    <div v-else class="flex justify-center items-center" style="height: 150px;">
      <p class="text-lg text-gray-400">No runs in last 24h.</p>
    </div>
  </div>
</template>
