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
                            <router-link :to="{ name: 'test-details', params: { id: run.testId } }">{{ tests[run.testId].name }}</router-link>
                        </td>
                        <td class="text-left">
                            <router-link :to="{ name: 'run-details', params: { id: run.id } }">view run</router-link>
                        </td>
                        <td class="text-right">{{ run.startedAt != null ? $filters.absoluteDate(run.startedAt) : null }}
                        </td>
                        <td class="text-right">
                            <span v-if="run.finishedAt != null && run.startedAt != null"> {{ $filters.relativeDate(run.startedAt, run.finishedAt, ['minutes', 'seconds']) }}</span> </td>
                        <td class="text-right">
                            <test-result-badge :result="run.result"></test-result-badge>
                        </td>
                        <td class="text-right">
                            <router-link :to="{ name: 'runner-details', params: { id: run.runnerId } }">{{ 'runnerId' in run ? runners[run.runnerId].name : '' }}</router-link>
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
    this.fetchRuns()
  },
  data() {
    return {
      runs: [],
      tests: {},
      runners: {}
    }
  },
  methods: {
    async fetchRuns() {
      const testIDs = new Set();
      const runnerIDs = new Set();
      const runs = await tstr.fetchRuns();
      runs.forEach(r => {
        testIDs.add(r.testId);
        if ('runnerId' in r) {
          runnerIDs.add(r.runnerId);
        }
      })

      const tests = await tstr.fetchTests(Array.from(testIDs.values()));
      const runners = await tstr.fetchRunners(Array.from(runnerIDs.values()));

      let testsByID = {}
      tests.forEach(t => {
        testsByID[t.id] = t
      })

      let runnersByID = {}
      await runners.forEach(r => {
        runnersByID[r.id] = r
      })

      this.runs = runs;
      this.tests = testsByID;
      this.runners = runnersByID;

      console.log("RUNS", this.runs);
      console.log("TESTS", this.tests);
      console.log("RUNNERS", this.runners);
    },
    
  }
}
</script>
