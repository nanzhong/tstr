import { createApp } from "vue";
import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";
import { createPinia } from "pinia";
import Notifications, { notify } from "notiwind";
import App from "./App.vue";
import "./index.css";

const Header = () => import("./components/Header.vue");
const NamespaceSelection = () => import("./views/NamespaceSelection.vue");
const Dashboard = () => import("./views/Dashboard.vue");
const RunDetails = () => import("./views/RunDetails.vue");
const RunnerDetails = () => import("./views/RunnerDetails.vue");
const Runners = () => import("./views/Runners.vue");
const Runs = () => import("./views/Runs.vue");
const TestDetails = () => import("./views/TestDetails.vue");
const Tests = () => import("./views/Tests.vue");

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(Notifications);
app.provide("apiPathPrefix", "/api");
app.config.globalProperties.$initReq = {
  pathPrefix: "/api",
};
app.config.errorHandler = (err, instance, info) => {
  console.log("[tstr ui error]", "err:", err, "instance:", instance, `in: ${info}`)
  notify({
    group: "top",
    type: "error",
    title: "Uh oh, something went wrong!",
    text: err.toString(),
  }, 4000);
};

function prefixRoutes(prefix: string, routes: RouteRecordRaw[]) {
  return routes.map((route) => {
    route.path = prefix + "/" + route.path;
    return route;
  })
}

const routes = [
  {
    path: "/",
    name: "home",
    redirect: "/nss",
  },
  {
    path: "/nss",
    name: "namespace-selection",
    props: {
      header: { title: "Select a Namespace" },
    },
    components: {
      default: NamespaceSelection,
      header: Header,
    },
  },
  ...prefixRoutes("/ns/:namespace", [
    {
      path: "dashboard",
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
      path: "tests",
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
      path: "runs",
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
      path: "runners",
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
      path: "tests/:id",
      name: "test-details",
      props: {
        header: { title: "Test Details" }
      },
      components: {
        default: TestDetails,
        header: Header,
      }
    },
    {
      path: "runs/:id",
      name: "run-details",
      props: {
        header: { title: "Run Details" }
      },
      components: {
        default: RunDetails,
        header: Header,
      }
    },
    {
      path: "runners/:id",
      name: "runner-details",
      props: {
        header: { title: "Runner Details" }
      },
      components: {
        default: RunnerDetails,
        header: Header,
      }
    },
  ]),
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

app.use(router);
app.mount("#app");
