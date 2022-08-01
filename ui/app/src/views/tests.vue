
<template>
    <q-page>
        <q-tab-panel name="tests">
            <div style="" v-for="test in tests" :key="test.id">
                <test-summary :id="test.id" :name="test.name"></test-summary>
            </div>
        </q-tab-panel>
    </q-page>
</template>

<script lang="ts" setup>
import TestSummary from '../components/TestSummary.vue'
</script>

<script lang="ts">
import { defineComponent } from 'vue';
import { DataService } from '../api/data/v1/data.pb';
import { Test } from '../api/common/v1/common.pb';


export default defineComponent({
    created() {
        (async () => {
            const tests = (await DataService.QueryTests({}, this.$initReq)).tests
            if (tests) {
                this.tests = tests
            }
        })()
    },
    data() {
        return {
            tests: [] as Test[],
        }
    },
})
</script>

<style>
</style>