<script setup lang="ts">
import dayjs from "dayjs";
import { useInitReq } from "../api/init";
import { DataService } from "../api/data/v1/data.pb";
import TestCard from "../components/TestCard.vue";

const initReq = useInitReq();

const tests = (await DataService.QueryTests({}, initReq)).tests!;
</script>

<template>
  <div class="max-w-7xl mx-auto pb-12 px-4 sm:px-6 lg:px-8">
    <div class="bg-white rounded-lg shadow px-5 py-6 sm:px-6">
      <div class="-mx-5 sm:-mx-6 px-5 sm:px-6 pb-5 border-b border-gray-200">
        <div class="-ml-2 -mt-2 flex flex-wrap items-baseline">
          <h2 class="ml-2 mt-2 text-lg leading-6 font-medium text-gray-900">All Tests</h2>
        </div>
      </div>

      <ul role="list" class="mt-5 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        <li v-for="test in tests" :key="test.id" class="col-span-1">
          <TestCard :test="test" />
        </li>
      </ul>
    </div>
  </div>
</template>
