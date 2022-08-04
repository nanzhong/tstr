<script setup lang="ts">
import { ref, onMounted } from "vue";
import dayjs from "dayjs";
import Plotly from "plotly.js-dist/plotly";
import { RunResult } from "../api/common/v1/common.pb";

const props = defineProps<{
  test: { id: string, name: string, results: { interval: dayjs.Dayjs, counts: { [key in RunResult]: number } }[] }
}>();

const plot = ref(null);

const hourBuckets: string[] = [];
const runSeries = {
  [RunResult.PASS]: { name: "Pass", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#22c55e" } },
  [RunResult.FAIL]: { name: "Fail", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ef4444" } },
  [RunResult.ERROR]: { name: "Error", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#ec4899" } },
  [RunResult.UNKNOWN]: { name: "Pending", type: "bar", x: hourBuckets, y: Array(0), marker: { color: "#0ea5e9" } },
}


props.test.results.forEach(r => {
  hourBuckets.push(r.interval.toISOString());
  for (const [result, count] of Object.entries(r.counts)) {
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
</script>

<template>
  <div class="bg-white rounded-lg shadow divide-y divide-gray-200">
    <router-link :to="{ name: 'test-details', params: { id: test.id } }">
      <h3 class="p-5 text-gray-900 text-sm font-medium truncate">{{ test.name }}</h3>
    </router-link>
    <div ref="plot"></div>
  </div>
</template>
