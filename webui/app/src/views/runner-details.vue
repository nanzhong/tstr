<script setup>
import TestResultBadge from '../components/TestResultBadge.vue'
import RunnerInfo from '../components/RunnerInfo.vue'
</script>

<template>

    <q-tab-panel name="runners" v-if="runner != null">
        <div class="text-h6">Runner: {{ runner.Name }}</div>

        <runner-info :runner="runner"></runner-info>

        <br />
        <div class="text-h6">Test Runs</div>

        <q-markup-table separator="vertical" flat bordered v-if="runner.LastRuns.length > 0">
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
                <tr v-for="run in runner.LastRuns">
                    <td class="text-left">
                        <router-link :to="{ name: 'test-details', params: { id: run.TestID } }"> {{ run.TestName }}
                        </router-link>
                    </td>
                    <td class="text-left">
                        <router-link :to="{ name: 'run-details', params: { id: run.RunID } }">view run</router-link>
                    </td>
                    <td class="text-right">{{ run.StartedAt != null ? $filters.absoluteDate(run.StartedAt) : null }}
                    </td>
                    <td class="text-right">
                        {{ $filters.relativeDate(run.StartedAt, run.FinishedAt, ['minutes', 'seconds']) }} </td>
                    <td class="text-right">
                        <test-result-badge :result="run.Result"></test-result-badge>
                    </td>
                </tr>

            </tbody>
        </q-markup-table>
        <div v-if="runner.LastRuns.length == 0">
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
        }
    },
    methods: {
        async fetchRunnerDetails(runnerId) {
            const runner = await tstr.fetchRunnerDetails(runnerId)
            this.runner = runner
        }
    }
}
</script>
