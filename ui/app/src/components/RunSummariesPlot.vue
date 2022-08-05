<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import Plotly from "plotly.js-dist/plotly";
import dayjs from "dayjs";
import duration from "dayjs/plugin/duration";
import relativeTime from "dayjs/plugin/relativeTime";
import { RunResult } from "../api/common/v1/common.pb";
import { RunSummary } from "../api/data/v1/data.pb";

dayjs.extend(duration)
dayjs.extend(relativeTime)

const router = useRouter();

const props = defineProps<{
  runSummaries: RunSummary[]
}>()

const plot = ref(null);

const series = {
  [RunResult.UNKNOWN]: { name: "Pending", type: "scatter", mode: "markers", x: Array(0), y: Array(0), hovertext: Array(0), summaries: Array(0), marker: { color: "#0ea5e9" } },
  [RunResult.FAIL]: { name: "Fail", type: "scatter", mode: "markers", x: Array(0), y: Array(0), hovertext: Array(0), summaries: Array(0), marker: { color: "#ef4444" } },
  [RunResult.ERROR]: { name: "Error", type: "scatter", mode: "markers", x: Array(0), y: Array(0), hovertext: Array(0), summaries: Array(0), marker: { color: "#ec4899" } },
  [RunResult.PASS]: { name: "Pass", type: "scatter", mode: "markers", x: Array(0), y: Array(0), hovertext: Array(0), summaries: Array(0), marker: { color: "#22c55e" } },
};

props.runSummaries.forEach(s => {
  let result = s.result;
  if (!result) { result = RunResult.UNKNOWN };
  series[result].summaries.push(s);
  series[result].x.push(dayjs(s.startedAt).toISOString());
  series[result].y.push(dayjs.duration(dayjs(s.finishedAt).diff(s.startedAt)).asSeconds());

  const items = [];
  items.push(`<b>${result}</b><br>`);
  if (s.resultData) {
    for (const [k, v] of Object.entries(s.resultData)) {
      items.push(`<b>${k}</b>: ${v}`);
    }
  }
  series[result].hovertext.push(items.join("<br>"));
});

const layout = {
  showlegend: true,
  xaxis: {
    type: 'date'
  },
  yaxis: {
    title: {
      text: "Duration (s)"
    }
  },
  height: 300,
  margin: {
    l: 40,
    r: 40,
    t: 40,
    b: 40,
  },
};

const onPlotlyClick = (data: any) => {
  const id = data.points[0].data.summaries[data.points[0].pointIndex].id;
  router.push({
    name: 'run-details',
    params: {
      id
    }
  })
};

onMounted(() => {
  Plotly.newPlot(plot.value, Object.values(series), layout, { responsive: true, displayModeBar: false });
  plot.value.on('plotly_click', onPlotlyClick);
});
</script>

<template>
  <div ref="plot"></div>
</template>
