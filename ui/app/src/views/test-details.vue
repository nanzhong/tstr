<script setup>
import TestDetails from '../components/TestDetails.vue'
import TestResultsChart from '../components/TestResultsChart.vue'
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

            <test-results-chart :runs="run_summaries"></test-results-chart>

        </q-tab-panel>
    </q-page>
</template>


<script>
import tstr from '../tstr'

export default {
    created() {
        this.fetchTestDetails(this.$route.params.id)
    },
    data() {
        return {
            test: null,
            run_summaries: null,
        }
    },
    methods: {
        async fetchTestDetails(testId) {
            const testDetails = await tstr.fetchTestDetails(testId)
            this.test = testDetails.test
            this.run_summaries = testDetails.run_summaries
        }
    }
}

</script>