<script setup lang="ts">
import { defineAsyncComponent } from 'vue';
const TestResultBadge = defineAsyncComponent(() => import('../components/TestResultBadge.vue'))
const RunnerInfo = defineAsyncComponent(() => import('../components/RunnerInfo.vue'))
const TestDetails = defineAsyncComponent(() => import('../components/TestDetails.vue'))
const TestLogLine = defineAsyncComponent(() => import('../components/TestLogLine.vue'))
const HumanDate = defineAsyncComponent(() => import('../components/HumanDate.vue'))
</script>

<template>
    <q-page>
        <q-inner-loading :showing="isLoading">
            <q-spinner-gears size="50px" color="primary" />
        </q-inner-loading>
        <q-tab-panel name="run" v-if="isLoading == false">

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
                                <span v-for="(v, k) in run.resultData"> {{ k }}: {{ v }}<br /></span><br />
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



<script lang="ts">

import tstr from '../tstr'
import { DataService } from '../api/data/v1/data.pb'
import { Run, Runner, Test } from '../api/common/v1/common.pb'

export default {
    created() {
        // this.fetchRunDetails(this.$route.params.id)
        (async () => {

            const run = await DataService.GetRun({ id: this.$route.params.id as string }, this.$initReq)

            const runner = await DataService.GetRunner({ id: run.run?.runnerId }, this.$initReq)

            const test = await DataService.GetTest({ id: run.run?.testId }, this.$initReq)

            this.run = run.run
            this.runner = runner.runner
            this.test = test.test

            console.log(this.run)

            this.isLoading = false
        })()
    },
    data() {
        return {
            run: undefined as (Run | undefined),
            runner: undefined as (Runner | undefined),
            test: undefined as (Test | undefined),
            isLoading: true,
        }
    },
}

</script>
