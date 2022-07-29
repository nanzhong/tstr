<script setup>
import TestResultBadge from '../components/TestResultBadge.vue'
import RunnerInfo from '../components/RunnerInfo.vue'
import TestDetails from '../components/TestDetails.vue'
import TestLogLine from '../components/TestLogLine.vue'
import HumanDate from '../components/HumanDate.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="run" v-if="runner != null && test != null">

            <div class="row inline">

                <div>
                    <div class="text-h6">Test: {{ test.name }}</div>
                    <test-details :test="test"></test-details>
                </div>
                <div class="">
                    <div class="text-h6">Runner: {{ runner.name }}</div>
                    <runner-info :runner="runner"></runner-info>
                </div>
                <div class="">
                    <div class="text-h6">Test Run</div>
                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Started</q-item-label>
                            <q-item-label>
                                <human-date :date="run.startedAt" :relative="false"></human-date>
                            </q-item-label>
                        </q-item-section>
                    </q-item>

                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Duration</q-item-label>
                            <q-item-label>
                                <human-date :date="run.startedAt" :diff="run.finishedAt" :relative="true"></human-date>
                            </q-item-label>
                        </q-item-section>
                    </q-item>

                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Result</q-item-label>
                            <q-item-label>
                                <test-result-badge :result="run.result"></test-result-badge>
                            </q-item-label>
                        </q-item-section>
                    </q-item>

                    <q-item>
                        <q-item-section>
                            <q-item-label caption>Result Data</q-item-label>
                            <q-item-label>
                                <span v-for="(v,k) in run.resultData"> {{ k }}: {{ v}}<br/></span><br/>
                            </q-item-label>
                        </q-item-section>
                    </q-item>
                </div>
            </div>

            <div class="text-h6">Logs</div>
            <test-log-line :logline="logline" v-for="logline in run.logs"></test-log-line>
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
        }
    },
    methods: {
        async fetchRunDetails(runId) {
            const runDetails = await tstr.fetchRunDetails(runId)
            this.run = runDetails
            console.log("RUN",this.run)

            const runnerDetails = await tstr.fetchRunnerDetails(this.run.runnerId, false)
            this.runner = runnerDetails.runner
            console.log("RUNNER",this.runner)

            const testDetails = await tstr.fetchTestDetails(this.run.testId, false)
            this.test = testDetails.test

            console.log("TEST",this.test)

        }
    }


}

</script>
