<style>
.result-square {
    width: 1em;
    height: 1em;
    display: inline-block;
    margin-right: 0.1em;
    margin-top: 0.1em;
    padding: 0;
}

.result-pass {
    background-color: green;
}

.result-unknown {
    background-color: black;
}

.result-fail {
    background-color: red;
}

.result-error {
    background-color: orange;
}
</style>

<template>
    <div class="text-h6">
        <router-link :to="{ name: 'test-details', params: { id: ID } }"> {{ Name }}
        </router-link>
    </div>
    <span style="" v-for="testRun in testDetails.RunsSummary"
        v-if="testDetails != null && testDetails.RunsSummary != null" class="result-square"
        :class="'result-square result-' + testRun.Result">
        &nbsp;
    </span>
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