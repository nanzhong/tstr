<script setup>
import TestDetails from '../components/TestDetails.vue'
import TestResultsChart from '../components/TestResultsChart.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="tests" v-if="test != null">
            <div class="row">
                <div>
                    <div class="text-h6">Test: {{ test.Name }}</div>
                    <test-details :test="test"></test-details>
                </div>
            </div>

            <test-results-chart :runs="test.RunsSummary"></test-results-chart>

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
        }
    },
    methods: {
        async fetchTestDetails(testId) {
            const testDetails = await tstr.fetchTestDetails(testId)
            this.test = testDetails
        }
    }
}

</script>