<script setup lang="ts">
import { inject } from "vue";
import dayjs from "dayjs";
import relativeTime from 'dayjs/plugin/relativeTime';
import { Run, Test } from "../api/common/v1/common.pb";
import { DataService, RunSummary } from "../api/data/v1/data.pb";
import RunResult from "./RunResult.vue";
import Labels from "./Labels.vue";
import { InitReq } from "../api/fetch.pb";
import TimeWithTooltip from "../components/TimeWithTooltip.vue";

dayjs.extend(relativeTime);

const props = defineProps<{
  runs: Run[] | RunSummary[]
}>();

const initReq: InitReq = inject('dataInitReq')!;
const testIDSet = new Set<string>();
props.runs.forEach(r => testIDSet.add(r.testId!));

const tests = (await DataService.QueryTests({ ids: Array.from(testIDSet) }, initReq)).tests!;
const testMap = new Map<string, Test>();
tests.forEach(t => testMap.set(t.id!, t));

</script>

<template>
  <table class="min-w-full divide-y divide-gray-300">
    <thead class="bg-gray-50">
      <tr>
        <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-500 sm:pl-6">Test</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-gray-500">Labels</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-gray-500">Result</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-gray-500">Scheduled</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-gray-500">Started</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold text-gray-500">Finished</th>
        <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
          <span class="sr-only">View</span>
        </th>
      </tr>
    </thead>
    <tbody class="divide-y divide-gray-200 bg-white">
      <tr v-for="run in runs" :key="run.id">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6"><router-link :to="{ name: 'test-details', params: { id: run.testId! } }">{{ testMap.get(run.testId!) ? testMap.get(run.testId!)!.name : run.testId }}</router-link></td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500"><Labels :labels="run.labels!" /></td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">
          <RunResult :result="run.result!" />
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">
          <TimeWithTooltip v-if="run.scheduledAt" :time="run.scheduledAt" :relative="true" />
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">
          <TimeWithTooltip v-if="run.startedAt" :time="run.startedAt" :relative="true" />
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">
          <TimeWithTooltip v-if="run.finishedAt" :time="run.finishedAt" :relative="true" />
        </td>
        <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
          <router-link :to="{ name: 'run-details', params: { id: run.id } }" custom v-slot="{ href, navigate }">
            <a :href="href" @click="navigate" class="text-indigo-600 hover:text-indigo-900">View<span
              class="sr-only">, {{ run.id }}</span></a>
          </router-link>
        </td>
      </tr>
    </tbody>
  </table>
</template>
