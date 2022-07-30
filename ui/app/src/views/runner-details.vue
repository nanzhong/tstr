<script lang="ts" setup>
import TestResultBadge from '../components/TestResultBadge.vue'
import RunnerInfo from '../components/RunnerInfo.vue'
import HumanDate from '../components/HumanDate.vue'
</script>

<template>

    <q-inner-loading :showing="isLoading">
        <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>
    <q-tab-panel name="runners" v-if="isLoading == false">
        <div class="text-h6">Runner: {{ runner?.name }}</div>

        <runner-info :runner="runner"></runner-info>

        <br />
        <div class="text-h6">Test Runs</div>

        <q-markup-table separator="vertical" flat bordered
            v-if="isLoading == false && typeof runSummaries !== 'undefined' && runSummaries != null && runSummaries.length > 0">
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
                <tr v-for="run in runSummaries">
                    <td class="text-left">
                        <router-link :to="{ name: 'test-details', params: { id: run.testId } }"> {{
                            testsByID.get(run.testId as string)?.name
                        }}
                        </router-link>
                    </td>
                    <td class="text-left">
                        <router-link :to="{ name: 'run-details', params: { id: run.id } }">view run</router-link>
                    </td>
                    <td class="text-right">
                        <human-date :relative="true" :date="run.startedAt"></human-date>
                    </td>
                    <td class="text-right">
                        <human-date :relative="true" :diff="run.startedAt" :date="run.finishedAt"></human-date>
                    </td>
                    <td class="text-right">
                        <test-result-badge :result="run.result"></test-result-badge>
                    </td>
                </tr>

            </tbody>
        </q-markup-table>
        <div v-if="typeof runSummaries === 'undefined' || runSummaries.length == 0">
            <p>
                This runner doesn't have any test run recorded.
            </p>
        </div>
    </q-tab-panel>
</template>

<script lang="ts">

import { DataService, RunSummary } from '../api/data/v1/data.pb';
import { InitReq } from '../api/fetch.pb';
import { Runner, Test } from '../api/common/v1/common.pb';
import { defineComponent } from 'vue';

const initReq: InitReq = {
    pathPrefix: '/api'
}

export default defineComponent({
    data() {
        return {
            isLoading: true,
            runner: undefined as (Runner | undefined),
            runSummaries: undefined as (RunSummary[] | undefined),
            testsByID: new Map<string,Test>(),
        }
    },
    created() {
        this.fetchRunnerDetails(this.$route.params.id as string)
    },
    methods: {
        async fetchRunnerDetails(runnerId: string) {
            const runnerDetails = await DataService.GetRunner({id: runnerId}, initReq)
            this.runner = runnerDetails.runner
            this.runSummaries = runnerDetails.runSummaries
            const testIDs: Set<string> = new Set(runnerDetails.runSummaries?.map( r => r.testId as string))
            const tests = await DataService.QueryTests({ids: Array.from(testIDs)},initReq)
            this.testsByID = (tests.tests?.filter(t => t) as Test[]).reduce((acc,test) => {
                acc.set(test.id as string, test)
                return acc
            }, new Map<string,Test>())

            this.isLoading = false
        }
    }
})

</script>
