<template>
    <q-badge color="red" v-if="result == 'fail'">{{ result }}</q-badge>
    <q-badge color="green" v-if="result == 'pass'">{{ result }}</q-badge>
    <q-badge color="black" v-if="result == 'unknown'">{{result}}</q-badge>
    <q-badge color="orange" v-if="result == 'error'">{{result}}</q-badge>
</template>

<script>
export default {
    props: ['result'],
};
</script>