<script setup lang="ts">
import { TestMatrixLabelValues } from "../api/common/v1/common.pb";

defineProps<{
  matrixLabels: { [key: string]: TestMatrixLabelValues } | undefined
}>();
</script>

<template>
  <table class="min-w-full divide-y divide-gray-300">
    <thead>
      <tr>
        <th scope="col" class="py-1.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6 md:pl-0">Label</th>
        <th scope="col" class="relative py-1.5 pl-3 text-left text-sm pr-4 sm:pr-6 md:pr-0">Values</th>
      </tr>
    </thead>
    <tbody class="divide-y divide-gray-200">
      <tr v-for="(values, label) in matrixLabels" :key="label">
        <td class="whitespace-nowrap py-1 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 md:pl-0">{{ label }}</td>
        <td class="relative whitespace-nowrap py-1 pl-3 pr-4 text-left text-sm font-medium sm:pr-6 md:pr-0">
          <span v-for="value in values.values" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500"> {{ label }}={{ value }} </span>
        </td>
      </tr>
    </tbody>
  </table>
</template>
