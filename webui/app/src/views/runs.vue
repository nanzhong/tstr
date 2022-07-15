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
                    <tr v-for="run in Runs">
                        <td class="text-left">
                            <router-link :to="{ name: 'test-details', params: { id: run.TestID } }"> {{ run.TestName }}
                            </router-link>
                        </td>
                        <td class="text-left">
                            <router-link :to="{ name: 'run-details', params: { id: run.ID } }">view run</router-link>
                        </td>
                        <td class="text-right">{{ run.StartedAt != null ? $filters.absoluteDate(run.StartedAt) : null }}
                        </td>
                        <td class="text-right">
                            {{ $filters.relativeDate(run.StartedAt, run.FinishedAt, ['minutes', 'seconds']) }} </td>
                        <td class="text-right">
                            <test-result-badge :result="run.Result"></test-result-badge>
                        </td>
                        <td class="text-right">
                            {{ run.RunnerName }}
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
            Runs: [],
        }
    },
    methods: {
        async fetchRunners() {
            const runs = await tstr.fetchRuns()
            this.Runs = runs
        }
    }

}
</script>