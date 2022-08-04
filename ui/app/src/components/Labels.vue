<script setup lang="ts">
defineProps<{
  labels: { [key: string]: string } | undefined
}>();
</script>

<template>
  <span v-for="(value, key) in labels" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500"> {{ key }}={{ value }} </span>
</template>
