<script setup lang="ts">
import { useRoute } from "vue-router";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { useInitReq } from "../api/init";
import { DataService } from "../api/data/v1/data.pb";
import Labels from "../components/Labels.vue";
import RunsTable from "../components/RunsTable.vue";
import TimeWithTooltip from "../components/TimeWithTooltip.vue";

dayjs.extend(relativeTime);

const route = useRoute();
const initReq = useInitReq();
const runnerDetails = await DataService.GetRunner({ id: route.params.id as string }, initReq);
const runner = runnerDetails.runner!;
const runSummaries = runnerDetails.runSummaries || [];
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">{{ runner.name }}</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">{{ runner. id }}</p>
        </div>
      </div>

      <div class="mt-5">
        <dl class="grid grid-cols-1 gap-x-2 gap-y-4 sm:grid-cols-2">
          <div class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Accept Label Selectors</dt>
            <dd class="mt-1 text-sm text-gray-900"><Labels :labels="runner.acceptTestLabelSelectors || null" /></dd>
          </div>
          <div class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Reject Label Selectors</dt>
            <dd class="mt-1 text-sm text-gray-900"><Labels :labels="runner.rejectTestLabelSelectors || null" /></dd>
          </div>
          <div class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Registered At</dt>
            <dd class="mt-1 text-sm text-gray-900"><TimeWithTooltip :time="runner.registeredAt" :relative="true" /></dd>
          </div>
          <div class="col-span-1">
            <dt class="text-sm font-medium text-gray-500">Last Heartbeat At</dt>
            <dd class="mt-1 text-sm text-gray-900"><TimeWithTooltip :time="runner.lastHeartbeatAt" :relative="true" /></dd>
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

      <div class="mt-5 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-6">
          <div class="inline-block min-w-full py-2 align-middle">
            <div class="overflow-hidden shadow-sm ring-1 ring-black ring-opacity-5">
              <RunsTable :runs="runSummaries" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
