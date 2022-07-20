<script setup>
import TestResultBadge from '../components/TestResultBadge.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="runs">
            <div class="text-h6">Latest runs</div>
            <q-markup-table separator="horizontal" flat bordered>
                <thead>
                    <tr>
                        <th class="text-left">Test</th>
                        <th class="text-left"></th>
                        <th class="text-right">Start time</th>
                        <th class="text-right">Duration</th>
                        <th class="text-right">Result</th>
                        <th class="text-right">Runner</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="run in runs">
                        <td class="text-left">
                            <router-link :to="{ name: 'test-details', params: { id: run.test_id } }"> FIXME test_name {{ run.TestName }}
                            </router-link>
                        </td>
                        <td class="text-left">
                            <router-link :to="{ name: 'run-details', params: { id: run.id } }">view run</router-link>
                        </td>
                        <td class="text-right">{{ run.started_at != null ? $filters.absoluteDate(run.started_at) : null }}
                        </td>
                        <td class="text-right">
                            <span v-if="run.finished_at != null && run.started_at != null"> {{ $filters.relativeDate(run.started_at, run.finished_at, ['minutes', 'seconds']) }}</span> </td>
                        <td class="text-right">
                            <test-result-badge :result="run.result"></test-result-badge>
                        </td>
                        <td class="text-right">
                            <router-link :to="{ name: 'runner-details', params: { id: run.runner_id } }">TODO runner_name {{run.runner_name}}</router-link>
                        </td>

                    </tr>

                </tbody>
            </q-markup-table>

        </q-tab-panel>
    </q-page>
</template>

<script>
import tstr from '../tstr'

export default {
    created() {
        this.fetchRunners()
    },
    data() {
        return {
            runs: [],
        }
    },
    methods: {
        async fetchRunners() {
            const runs = await tstr.fetchRuns()
            this.runs = runs
            console.log("RUNS",this.runs)
        }
    }

}
</script>