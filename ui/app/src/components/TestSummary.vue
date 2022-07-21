<script setup>
import TestResultsChart from '../components/TestResultsChart.vue'
</script>


<template>
    <div class="text-h6">
        <router-link :to="{ name: 'test-details', params: { id: id } }"> {{ name }}
        </router-link>
    </div>
    <test-results-chart :runs="runSummaries" v-if="runSummaries != null"></test-results-chart>
    <div style="margin-bottom: 2em"></div>
</template>

<script>
import tstr from '../tstr'

export default {
    created() {
        this.fetchTestDetails()
    },
    data() {
        return {
            test: null,
            runSummaries: null
        }
    },
    props: {
        id: { required: true, type: String },
        name: { required: true, type: String },
    },
    methods: {
        async fetchTestDetails() {
            const testDetails = await tstr.fetchTestDetails(this.id)
            this.test = testDetails.test
            this.runSummaries = testDetails.runSummaries
            console.log(this.runSummaries)
        }
    }

};
</script>
