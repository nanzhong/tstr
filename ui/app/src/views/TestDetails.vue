<script setup lang="ts">
import { useRoute } from "vue-router";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { useInitReq } from "../api/init";
import { DataService } from "../api/data/v1/data.pb";
import RunSummariesPlot from "../components/RunSummariesPlot.vue";
import Labels from "../components/Labels.vue";
import MatrixLabels from "../components/MatrixLabels.vue";

dayjs.extend(relativeTime);

const initReq = useInitReq();
const route = useRoute();
const testData = await DataService.GetTest({ id: route.params.id as string }, initReq);
const test = testData.test!;
const runSummaries = testData.runSummaries!;
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">{{ test.name }}</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">{{ test. id }}</p>
        </div>
      </div>

      <div class="mt-5">
        <dl class="grid grid-cols-1 gap-x-2 gap-y-4 mb-4 sm:grid-cols-3">
          <div class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Container Image</dt>
            <dd class="mt-1 text-sm text-gray-900"><pre>{{ test.runConfig?.containerImage }}</pre></dd>
          </div>
          <div v-if="test.runConfig?.command" class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Command</dt>
            <dd class="mt-1 text-sm text-gray-900"><pre>{{ test.runConfig?.command }}</pre></dd>
          </div>
          <div v-if="test.runConfig?.args"  class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Args</dt>
            <dd class="mt-1 text-sm text-gray-900">
              <pre><span v-for="(arg, i) in test.runConfig?.args">"{{ arg }}"<span v-if="i < test.runConfig?.args.length - 1">, </span></span></pre></dd>
          </div>
          <div v-if="test.cronSchedule" class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Cron Schedule</dt>
            <dd class="mt-1 text-sm text-gray-900"><pre>{{ test.cronSchedule }}</pre></dd>
          </div>
          <div lass="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Next Run At</dt>
            <dd class="mt-1 text-sm text-gray-900"><pre>{{ dayjs(test.nextRunAt!).fromNow() }}</pre></dd>
          </div>
          <div v-if="test.runConfig?.timeout" class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Timeout</dt>
            <dd class="mt-1 text-sm text-gray-900"><pre>{{ test.runConfig?.timeout }}</pre></dd>
          </div>
        </dl>
        <dl class="grid grid-cols-1 gap-x-2 gap-y-4 mt-4 sm:grid-cols-3">
          <div :class="[test.matrix?.labels ? 'sm:col-span-1' : 'sm:col-span-3', 'col-span-1']">
            <dt class="text-sm font-medium text-gray-500">Labels</dt>
            <dd class="mt-1 text-sm text-gray-900"><Labels :labels="test.labels" /></dd>
          </div>
          <div v-if="test.matrix?.labels" class="col-span-1 sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500">Matrix Labels</dt>
            <dd class="mt-1 text-sm text-gray-900"><MatrixLabels :matrixLabels="test.matrix?.labels!" /></dd>
          </div>
        </dl>
      </div>
    </div>

    <div class="mt-5 bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Runs</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <RunSummariesPlot v-if="runSummaries" :runSummaries="runSummaries" />
    <div v-else>
      <p class="mt-5 text-gray-400">No runs in last 24h.</p>
    </div>
    </div>
  </div>
</template>
