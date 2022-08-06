<script setup lang="ts">
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);

const props = defineProps<{
  time: any,
  relative: boolean,
}>();

const parsedTime = dayjs(props.time);

</script>

<template>
  <span class="has-tooltip">
    <span class="tooltip rounded shadow-lg p-2 bg-gray-500 text-white -mt-10">{{ relative ? parsedTime.toISOString() : parsedTime.fromNow() }}</span>
    <span>{{ relative ? parsedTime.fromNow() : parsedTime.toISOString() }}</span>
  </span>
</template>
