<script setup>
import TestResultsChart from '../components/TestResultsChart.vue'
</script>


<template>
    <div class="text-h6">
        <router-link :to="{ name: 'test-details', params: { id: ID } }"> {{ Name }}
        </router-link>
    </div>
    <test-results-chart :runs="testDetails.RunsSummary" v-if=" testDetails != null && testDetails.RunsSummary != null"></test-results-chart>
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
            testDetails: null
        }
    },
    props: {
        ID: { required: true, type: String },
        Name: { required: true, type: String },
    },
    methods: {
        async fetchTestDetails() {
            const testDetails = await tstr.fetchTestDetails(this.ID)
            this.testDetails = testDetails;
        }
    }

};
</script>