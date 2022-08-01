<script setup lang="ts">
import { defineAsyncComponent } from 'vue'

const TestDetails = defineAsyncComponent(() => import('../components/TestDetails.vue'))
const TestResultsChart = defineAsyncComponent(() => import('../components/TestResultsChart.vue'))
</script>

<template>
    <q-page>
        <q-tab-panel name="tests" v-if="test != null">
            <div class="row">
                <div>
                    <div class="text-h6">Test: {{ test.name }}</div>
                    <test-details :test="test"></test-details>
                </div>
            </div>

            <test-results-chart :runs="runSummaries"></test-results-chart>

        </q-tab-panel>
    </q-page>
</template>


<script lang="ts">

import { Test } from '../api/common/v1/common.pb';
import { DataService, RunSummary } from '../api/data/v1/data.pb';

export default {
    created() {
        (async () => {
            const testDetails = (await DataService.GetTest({ id: this.$route.params.id as string }, this.$initReq))
            if (testDetails.test) {
                this.test = testDetails.test
                console.log(this.test)
            }

            if (testDetails.runSummaries) {
                this.runSummaries = testDetails.runSummaries
            }
        })()

    },
    data() {
        return {
            test: undefined as Test | undefined,
            runSummaries: [] as RunSummary[],
        }
    },
}

</script>
