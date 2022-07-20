<script setup>
import TestResultsChart from '../components/TestResultsChart.vue'
</script>


<template>
    <div class="text-h6">
        <router-link :to="{ name: 'test-details', params: { id: id } }"> {{ name }}
        </router-link>
    </div>
    <test-results-chart :runs="run_summaries" v-if="run_summaries != null"></test-results-chart>
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
            run_summaries: null
        }
    },
    props: {
        id: { required: true, type: String },
        name: { required: true, type: String },
    },
    methods: {
        async fetchTestDetails() {
            const test_details = await tstr.fetchTestDetails(this.id)
            this.test = test_details.test
            this.run_summaries = test_details.run_summaries
            console.log(this.run_summaries)
        }
    }

};
</script>