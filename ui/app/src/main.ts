import { createApp, defineAsyncComponent, VueElement } from "vue";
import { Quasar } from "quasar";
import { createRouter, createWebHashHistory } from "vue-router";

import "@quasar/extras/material-icons/material-icons.css";
import "quasar/src/css/index.sass";

const App = defineAsyncComponent(() => import("./App.vue"));
const RunDetails = defineAsyncComponent( () => import("./views/run-details.vue"));
const RunnerDetails = defineAsyncComponent( () => import("./views/runner-details.vue"));
const Runners = defineAsyncComponent(() => import("./views/runners.vue"));
const Runs = defineAsyncComponent(() => import("./views/runs.vue"));
const TestDetails = defineAsyncComponent( () => import("./views/test-details.vue"));
const Tests = defineAsyncComponent(() => import("./views/tests.vue"));
const app = createApp(App);

app.config.globalProperties.$initReq = {
  pathPrefix: "/api",
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

app.use(router);

app.use(Quasar, {
  plugins: {},
});

app.mount("#app");
