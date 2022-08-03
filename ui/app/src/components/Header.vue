<script setup lang="ts">
defineProps([ "title", "subtitle"]);
</script>

<template>
  <h1 class="text-3xl font-bold text-white">{{ title }}<span v-if="subtitle !== ''">{{ subtitle }}</span></h1>
</template>
