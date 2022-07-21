<template>
    <q-badge color="red" v-if="result == 'FAIL'">{{ result }}</q-badge>
    <q-badge color="green" v-if="result == 'PASS'">{{ result }}</q-badge>
    <q-badge color="black" v-if="result == 'UNKNOWN'">{{result}}</q-badge>
    <q-badge color="orange" v-if="result == 'ERROR'">{{result}}</q-badge>
</template>

<script>
export default {
    props: ['result'],
};
</script>