import { createApp, defineAsyncComponent } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import App from "./App.vue";
import "./index.css";

const Header = () => import("./components/Header.vue");
const Dashboard = () => import("./views/Dashboard.vue");
const RunDetails = () => import("./views/RunDetails.vue");
const RunnerDetails = () => import("./views/RunnerDetails.vue");
const Runners = () => import("./views/Runners.vue");
const Runs = () => import("./views/Runs.vue");
const TestDetails = () => import("./views/TestDetails.vue");
const Tests = () => import("./views/Tests.vue");

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
    props: {
      header: { title: "Tests" },
    },
    components: {
      default: Tests,
      header: Header,
    }
  },
  {
    path: "/runs",
    name: "runs",
    props: {
      header: { title: "Runs" }
    },
    components: {
      default: Runs,
      header: Header,
    }
  },
  {
    path: "/runners",
    name: "runners",
    props: {
      header: { title: "Runners" }
    },
    components: {
      default: Runners,
      header: Header,
    }
  },
  {
    name: "test-details",
    path: "/tests/:id",
    props: {
      header: { title: "Test Details" }
    },
    components: {
      default: TestDetails,
      header: Header,
    }
  },
  {
    name: "run-details",
    path: "/runs/:id",
    props: {
      header: { title: "Run Details" }
    },
    components: {
      default: RunDetails,
      header: Header,
    }
  },
  {
    name: "runner-details",
    path: "/runners/:id",
    props: {
      header: { title: "Runner Details" }
    },
    components: {
      default: RunnerDetails,
      header: Header,
    }
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

app.use(router);
app.mount("#app");
