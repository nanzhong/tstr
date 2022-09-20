<script setup lang="ts">
import { useRoute } from "vue-router";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import duration from "dayjs/plugin/duration";
import { useInitReq } from "../api/init";
import { DataService } from "../api/data/v1/data.pb";
import { RunLogOutput, Runner } from "../api/common/v1/common.pb";
import RunResult from "../components/RunResult.vue";
import Labels from "../components/Labels.vue";
import TimeWithTooltip from "../components/TimeWithTooltip.vue";

dayjs.extend(relativeTime);
dayjs.extend(duration)

const route = useRoute();
const initReq = useInitReq();
const run = (await DataService.GetRun({ id: route.params.id as string }, initReq)).run!;
const test = (await DataService.GetTest({ id: run.testId }, initReq)).test!;
let runner: Runner | undefined;
// TODO this is a quirk where the zero value of a uuid is not the same as the
// zero value of a string pb field. We should fix this at the api level.
if (run.runnerId !== "00000000-0000-0000-0000-000000000000") {
 runner = (await DataService.GetRunner({ id: run.runnerId }, initReq)).runner;
}

const logColour = (t: RunLogOutput | null) => {
  switch (t) {
    case RunLogOutput.STDOUT:
      return "border-green-300";
    case RunLogOutput.STDOUT:
      return "border-red-300";
    case RunLogOutput.TSTR:
      return "border-pink-300";
    default:
      return "border-gray-300";
  }
};

const logDecode = (data: string | null) => {
  if (!data || data.length === 0) {
    return "";
  }
  return atob(data);
};
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
      <div class="col-span-1 md:col-span-3 bg-white rounded-lg shadow px-5 py-6 sm:px-6">
        <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
          <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
            <RunResult class="text-lg" :result="run.result!" />
            <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Run for {{ test.name }}</h2>
            <p class="ml-2 mt-1 text-sm text-gray-500 truncate">{{ run.id }}</p>
          </div>
        </div>

        <dl class="grid grid-cols-1 gap-2 mt-4 sm:grid-cols-3 text-sm">
          <div class="col-span-1">
            <dt class="font-medium text-gray-500">Scheduled At</dt>
            <dd class="mt-1 text-gray-900">
              <TimeWithTooltip :time="run.scheduledAt" :relative="true" />
            </dd>
          </div>
          <div class="col-span-1">
            <dt class="font-medium text-gray-500">Started At</dt>
            <dd class="mt-1 text-gray-900">
              <TimeWithTooltip v-if="run.startedAt" :time="run.startedAt" :relative="true" />
            </dd>
          </div>
          <div class="col-span-1">
            <dt class="font-medium text-gray-500">Finished At</dt>
            <dd class="mt-1 text-gray-900">
              <TimeWithTooltip v-if="run.finishedAt" :time="run.finishedAt" :relative="true" />
              <span v-if="run.finishedAt"> ({{ dayjs.duration(dayjs(run.finishedAt).diff(run.startedAt)).humanize() }})</span>
            </dd>
          </div>
        </dl>
      </div>

      <div class="col-span-1">
        <div class="bg-white rounded-lg shadow px-5 py-6 pb-10 sm:px-6 relative overflow-hidden">
          <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
            <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
              <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Test</h2>
              <p class="ml-2 mt-1 text-xs text-gray-500">This is the configuration for this run. It may differ from the current active configuration of the test.</p>
            </div>
          </div>

          <dl class="grid grid-cols-1 gap-2 my-4">
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500">Container Image</dt>
              <dd class="mt-1 text-sm text-gray-900"><pre>{{ run.testRunConfig?.containerImage }}</pre></dd>
            </div>
            <div v-if="run.testRunConfig?.command" class="col-span-1">
              <dt class="text-sm font-medium text-gray-500">Command</dt>
              <dd class="mt-1 text-sm text-gray-900"><pre>{{ run.testRunConfig?.command }}</pre></dd>
            </div>
            <div v-if="run.testRunConfig?.args" class="col-span-1">
              <dt class="text-sm font-medium text-gray-500">Args</dt>
              <dd class="mt-1 text-sm text-gray-900">
                <pre><span v-for="(arg, i) in run.testRunConfig?.args">"{{ arg }}"<span v-if="i < run.testRunConfig?.args.length - 1">, </span></span></pre></dd>
            </div>
            <div v-if="run.testRunConfig?.timeout" class="col-span-1">
              <dt class="text-sm font-medium text-gray-500">Timeout</dt>
              <dd class="mt-1 text-sm text-gray-900"><pre>{{ run.testRunConfig?.timeout }}</pre></dd>
            </div>
          </dl>

          <div class="absolute bottom-0 inset-x-0 bg-gray-50 px-4 py-2 sm:px-6">
            <div class="text-sm">
              <router-link :to="{ name: 'test-details', params: { id: test.id } }" class="font-medium text-indigo-600 hover:text-indigo-500">View<span class="sr-only"> {{ test.name }} details</span></router-link>
            </div>
          </div>
        </div>

        <div v-if="runner" class="bg-white rounded-lg shadow px-5 py-6 pb-10 sm:px-6 mt-5 relative overflow-hidden">
          <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
            <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
              <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Runner</h2>
            </div>
          </div>

          <dl class="grid grid-cols-1 gap-2 my-4 text-sm">
            <div class="col-span-1">
              <dt class="font-medium text-gray-500">Name</dt>
              <dd class="mt-1 text-gray-900">{{ runner.name }}</dd>
            </div>
            <div class="col-span-1" v-if="runner.acceptTestLabelSelectors">
              <dt class="font-medium text-gray-500">Accept Label Selectors</dt>
              <dd class="mt-1 text-gray-900"><Labels :labels="runner.acceptTestLabelSelectors" /></dd>
            </div>
            <div class="col-span-1" v-if="runner.rejectTestLabelSelectors">
              <dt class="font-medium text-gray-500">Reject Label Selectors</dt>
              <dd class="mt-1 text-gray-900"><Labels :labels="runner.rejectTestLabelSelectors" /></dd>
            </div>
            <div class="col-span-1">
              <dt class="font-medium text-gray-500">Last Heartbeat At</dt>
              <dd class="mt-1 text-gray-900"><TimeWithTooltip :time="runner.lastHeartbeatAt" :relative="true" /></dd>
            </div>
          </dl>

          <div class="absolute bottom-0 inset-x-0 bg-gray-50 px-4 py-2 sm:px-6">
            <div class="text-sm">
              <router-link :to="{ name: 'runner-details', params: { id: runner.id } }" class="font-medium text-indigo-600 hover:text-indigo-500">View<span class="sr-only"> {{ runner.name }} details</span></router-link>
            </div>
          </div>
        </div>
      </div>

      <div class="col-span-1 md:col-span-2">
        <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
          <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
            <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
              <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Logs</h2>
            </div>
          </div>

          <div class="mt-5 text-xs overflow-scroll">
            <div v-for="log in run.logs" :class="[logColour(log.outputType), 'border-l-4 pl-2']">
              <pre><span v-if="log.time" class="text-gray-400">{{ log.time }} </span>{{ logDecode(log.data) }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!--
  <q-page>
      <div class="row inline">

        <div>
          <div class="text-h6">Test: {{ test.name }}</div>
          <test-details :test="test"></test-details>
        </div>
        <div class="">
          <div class="text-h6">Runner: {{ runner.name }}</div>
          <runner-info :runner="runner"></runner-info>
        </div>
        <div class="">
          <div class="text-h6">Test Run</div>
          <q-item>
            <q-item-section>
              <q-item-label caption>Started</q-item-label>
              <q-item-label>
                <human-date :date="run.startedAt" :relative="false"></human-date>
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section>
              <q-item-label caption>Duration</q-item-label>
              <q-item-label>
                <human-date :date="run.startedAt" :diff="run.finishedAt" :relative="true"></human-date>
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section>
              <q-item-label caption>Result</q-item-label>
              <q-item-label>
                <test-result-badge :result="run.result"></test-result-badge>
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item>
            <q-item-section>
              <q-item-label caption>Result Data</q-item-label>
              <q-item-label>
                <span v-for="(v, k) in run.resultData"> {{ k }}: {{ v }}<br /></span><br />
              </q-item-label>
            </q-item-section>
          </q-item>
        </div>
      </div>

      <div class="text-h6">Logs</div>
      <test-log-line :logline="logline" v-for="logline in run.logs"></test-log-line>
  </q-page>
    -->
</template>
