import { createApp, defineAsyncComponent } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import { createPinia, storeToRefs } from "pinia";
import { useNamespaceStore } from "./stores/namespace";
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
app.provide("apiPathPrefix", "/api");
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
    path: "/namespaces",
    name: "namespace-selection",
    props: {
      header: { title: "Select a Namespace" },
    },
    components: {
      default: NamespaceSelection,
      header: Header,
    },
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

router.beforeEach((to, from) => {
  const nsStore = useNamespaceStore();
  const { currentNamespace } = storeToRefs(nsStore);
  if (to.name !== "namespace-selection" && currentNamespace.value === "") {
    return { name: "namespace-selection" };
  }
});

app.use(router);
app.mount("#app");
