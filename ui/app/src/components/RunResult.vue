<script setup lang="ts">
import { RunResult } from "../api/common/v1/common.pb";
defineProps<{
  result: RunResult
}>();
</script>

<template>
  <span v-if="result === RunResult.FAIL" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-500"> Fail </span>
  <span v-else-if="result === RunResult.PASS" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-500"> Pass </span>
  <span v-else-if="result === RunResult.UNKNOWN" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-sky-500"> Pending </span>
  <span v-else-if="result === RunResult.ERROR" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-pink-100 text-pink-500"> Error </span>
</template>
