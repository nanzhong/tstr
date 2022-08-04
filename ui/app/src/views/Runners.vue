<script setup lang="ts">
import { inject } from "vue";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { InitReq } from "../api/fetch.pb";
import { DataService } from "../api/data/v1/data.pb";
import Labels from "../components/Labels.vue";

dayjs.extend(relativeTime);

const initReq: InitReq = inject('dataInitReq')!;

const runners = (await DataService.QueryRunners({ lastHeartbeatWithin: `${24*60*60}s` }, initReq)).runners!;
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">Active Runners</h2>
          <p class="ml-2 mt-1 text-sm text-gray-500 truncate">last 24h</p>
        </div>
      </div>

      <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-6">
        <div class="inline-block min-w-full py-2 align-middle">
          <div class="overflow-hidden shadow-sm ring-1 ring-black ring-opacity-5">
            <table class="min-w-full divide-y divide-gray-300">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold uppercase text-gray-500 sm:pl-6">Name</th>
                  <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Accept Label Selectors</th>
                  <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Reject Label Selector</th>
                  <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Registered At</th>
                  <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Last Heartbeat at</th>
                  <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                    <span class="sr-only">View</span>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 bg-white">
                <tr v-for="runner in runners" :key="runner.id">
                  <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">{{ runner.name }}</td>
                  <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500"><Labels :labels="runner.acceptTestLabelSelectors || null" /></td>
                  <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500"><Labels :labels="runner.rejectTestLabelSelectors || null" /></td>
                  <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">{{ dayjs(runner.registeredAt).fromNow() }}</td>
                  <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500">{{ dayjs(runner.lastHeartbeatAt).fromNow() }}</td>
                  <td class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                    <router-link :to="{ name: 'runner-details', params: { id: runner.id } }" custom v-slot="{ href, navigate }">
                      <a :href="href" @click="navigate" class="text-indigo-600 hover:text-indigo-900">View<span
                        class="sr-only">, {{ runner.name }}</span></a>
                    </router-link>
                  </td>
                </tr>
              </tbody>
            </table>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>
