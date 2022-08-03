import { createApp, defineAsyncComponent } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";

import "./index.css";

const App = defineAsyncComponent(() => import("./App.vue"));
const Header = () => import("./components/Header.vue");
const Dashboard = () => import("./views/Dashboard.vue");
const RunDetails = () => import("./views/run-details.vue");
const RunnerDetails = () => import("./views/runner-details.vue");
const Runners = () => import("./views/runners.vue");
const Runs = () => import("./views/runs.vue");
const TestDetails = () => import("./views/test-details.vue");
const Tests = () => import("./views/tests.vue");

const app = createApp(App);
app.provide("dataInitReq", { pathPrefix: "/api" });
app.config.globalProperties.$initReq = {
  pathPrefix: "/api",
};

const routes = [
  {
    path: "/",
    name: "home",
    redirect: "/dashboard",
  },
  {
    path: "/dashboard",
    name: "dashboard",
    props: { 
      header: { title: "Dashboard" },
    },
    components: {
      default: Dashboard,
      header: Header,
    },
  },
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
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

app.use(router);
app.mount("#app");
