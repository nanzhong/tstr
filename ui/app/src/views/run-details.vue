<script setup>
import TestResultBadge from '../components/TestResultBadge.vue'
import RunnerInfo from '../components/RunnerInfo.vue'
import TestDetails from '../components/TestDetails.vue'
import TestLogLine from '../components/TestLogLine.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="run" v-if="runner != null && test != null">

            <div class="row inline">

                <div>
                    <div class="text-h6">Test: {{ test.Name }}</div>
                    <test-details :test="test"></test-details>
                </div>
                <div class="">
                    <div class="text-h6">Runner: {{ runner.Name }}</div>
                    <runner-info :runner="runner"></runner-info>
                </div>
                <div class="">
                    <div class="text-h6">Test Run</div>
                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Started</q-item-label>
                            <q-item-label>
                                {{ $filters.absoluteDate(run.StartedAt) }}
                            </q-item-label>
                        </q-item-section>
                    </q-item>

                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Duration</q-item-label>
                            <q-item-label>
                                {{ $filters.relativeDate(run.StartedAt, run.FinishedAt) }}
                            </q-item-label>
                        </q-item-section>
                    </q-item>

                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Result</q-item-label>
                            <q-item-label>
                                <test-result-badge :result="run.Result"></test-result-badge>
                            </q-item-label>
                        </q-item-section>
                    </q-item>
                </div>
            </div>

            <div class="text-h6">Logs</div>
            <test-log-line :logline="logline" v-for="logline in run.Logs"></test-log-line>
        </q-tab-panel>
    </q-page>


</template>



<script>

import tstr from '../tstr'

export default {
    created() {
        this.fetchRunDetails(this.$route.params.id)
    },
    data() {
        return {
            run: null,
            runner: null,
            test: null,
            runconfig: null,
        }
    },
    methods: {
        async fetchRunDetails(runId) {
            const runDetails = await tstr.fetchRunDetails(runId)
            this.run = runDetails

            const runnerDetails = await tstr.fetchRunnerDetails(runDetails.RunnerID, false)
            this.runner = runnerDetails

            const testDetails = await tstr.fetchTestDetails(runDetails.TestID, false)
            this.test = testDetails
        }
    }


}

</script>