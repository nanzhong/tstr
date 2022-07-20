<script setup>
import TestResultBadge from '../components/TestResultBadge.vue'
import RunnerInfo from '../components/RunnerInfo.vue'
</script>

<template>

    <q-tab-panel name="runners" v-if="runner != null">
        <div class="text-h6">Runner: {{ runner.name }}</div>

        <runner-info :runner="runner"></runner-info>

        <br />
        <div class="text-h6">Test Runs</div>

        <q-markup-table separator="vertical" flat bordered v-if="run_summaries != null && run_summaries.length > 0">
            <thead>
                <tr>
                    <th class="text-left">Test</th>
                    <th class="text-left"></th>
                    <th class="text-right">Start time</th>
                    <th class="text-right">Duration</th>
                    <th class="text-right">Result</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="run in run_summaries">
                    <td class="text-left">
                        <router-link :to="{ name: 'test-details', params: { id: run.test_id } }"> TODO:test_name {{ run.test_name }}
                        </router-link>
                    </td>
                    <td class="text-left">
                        <router-link :to="{ name: 'run-details', params: { id: run.id } }">view run</router-link>
                    </td>
                    <td class="text-right">{{ run.started_at != null ? $filters.absoluteDate(run.started_at) : null }}
                    </td>
                    <td class="text-right">
                        {{ $filters.relativeDate(run.started_at, run.finished_at, ['minutes', 'seconds']) }} </td>
                    <td class="text-right">
                        <test-result-badge :result="run.result"></test-result-badge>
                    </td>
                </tr>

            </tbody>
        </q-markup-table>
        <div v-if="run_summaries.length == 0">
            <p>
                This runner doesn't have any test run recorded.
            </p>
        </div>
    </q-tab-panel>
</template>

<script>

import tstr from '../tstr'

export default {
    created() {
        this.fetchRunnerDetails(this.$route.params.id)
    },
    data() {
        return {
            runner: null,
            run_summaries: null,
        }
    },
    methods: {
        async fetchRunnerDetails(runnerId) {
            const runnerDetails = await tstr.fetchRunnerDetails(runnerId)
            this.runner = runnerDetails.runner
            this.run_summaries = runnerDetails.run_summaries
            console.log("RUNNER", this.runner)
            console.log("RUN_SUMMARIES", this.run_summaries)
        }
    }
}
</script>
