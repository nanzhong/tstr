import { createApp } from "vue";
import { Quasar } from "quasar";
import { createRouter, createWebHashHistory } from "vue-router";

import "@quasar/extras/material-icons/material-icons.css";
import "quasar/src/css/index.sass";

import VueApexCharts from "vue3-apexcharts";

// import './style.css'
import App from "./App.vue";

import RunDetails from "./views/run-details.vue";
import RunnerDetails from "./views/runner-details.vue";
import Runners from "./views/runners.vue";
import Runs from "./views/runs.vue";
import TestDetails from "./views/test-details.vue";
import Tests from "./views/tests.vue";
import { DateTime, Interval, Duration } from "luxon";

const app = createApp(App);

app.config.globalProperties.$filters = {
  relativeDate(
    date,
    date2 = null,
    diffUnits = ["months", "days", "hours", "minutes", "seconds"]
  ) {
    if (date == null) {
      return null;
    }
    var i;
    if (date2 == null) date2 = DateTime.now();
    if (date2 > date) {
      i = Interval.fromDateTimes(date, date2);
    } else {
      i = Interval.fromDateTimes(date2, date);
    }
    return i.toDuration(diffUnits).toHuman({ unitDisplay: "short" });
  },
  absoluteDate(date) {
    if (date == null) return null;

    return date.toHTTP();
  },
};

const routes = [
  {
    path: "/tests",
    name: "tests",
    component: Tests,
  },
  {
    path: "/runs",
    name: "runs",
    component: Runs,
  },
  {
    path: "/runners",
    name: "runners",
    component: Runners,
  },
  {
    name: "test-details",
    path: "/tests/:id",
    component: TestDetails,
  },
  {
    name: "run-details",
    path: "/runs/:id",
    component: RunDetails,
  },
  {
    name: "runner-details",
    path: "/runners/:id",
    component: RunnerDetails,
  },
  {
    path: "/",
    name: "home",
    component: Tests,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

app.use(VueApexCharts);

app.use(router);

app.use(Quasar, {
  plugins: {},
});

app.mount("#app");
