<script setup lang="ts">

const props = defineProps<{
  text: any
}>();

</script>

<template>
  <span class="has-tooltip relative">
    <span class="tooltip rounded shadow-lg p-2 bg-gray-500 text-white -mt-10">{{ text }}</span>
    <pre class="truncate">{{ text }}</pre>
  </span>
</template>
