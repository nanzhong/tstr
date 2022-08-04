<script setup lang="ts">
import { RunResult } from "../api/common/v1/common.pb";

const props = defineProps<{
  testRunCounts: { testID: string, testName: string, runCounts: { [key in RunResult]: number } }[]
}>();
</script>

<template>
  <table class="min-w-full divide-y divide-gray-300">
    <thead class="bg-gray-50">
      <tr>
        <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold uppercase text-gray-500 sm:pl-6">Test</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Scheduled</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Pass</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Fail</th>
        <th scope="col" class="px-2 py-3.5 text-left text-sm font-semibold uppercase text-gray-500">Error</th>
        <th scope="col" class="py-3.5 pl-3 pr-4 text-left text-sm font-semibold uppercase text-gray-500 sm:pr-6">Unknown</th>
      </tr>
    </thead>
    <tbody class="divide-y divide-gray-200 bg-white">
      <tr v-for="count in testRunCounts" :key="count.testID" :set="totalRunCount=count.runCounts[RunResult.PASS] + count.runCounts[RunResult.FAIL] + count.runCounts[RunResult.ERROR] + count.runCounts[RunResult.UNKNOWN]">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6"><router-link :to="{ name: 'test-details', params: { id: count.testID! } }">{{ count.testName }}</router-link></td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-500 tabular-nums">{{ totalRunCount }}</td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-green-500 tabular-nums">{{ count.runCounts[RunResult.PASS] }}</td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-red-500 tabular-nums">{{ count.runCounts[RunResult.FAIL] }}</td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-pink-500 tabular-nums">{{ count.runCounts[RunResult.ERROR] }}</td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm text-gray-500 tabular-nums sm:pr-6">{{ count.runCounts[RunResult.UNKNOWN] }}</td>
      </tr>
    </tbody>
  </table>
</template>
