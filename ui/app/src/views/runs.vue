<script setup>
import { defineAsyncComponent } from 'vue'
const HumanDate = defineAsyncComponent(() => import('../components/HumanDate.vue'))
const TestResultBadge = defineAsyncComponent(() => import('../components/TestResultBadge.vue'))
</script>

<template>
  <q-page>
    <q-inner-loading :showing="isLoading">
        <q-spinner-gears size="50px" color="primary" />
    </q-inner-loading>
    <q-tab-panel name="runs" v-if="isLoading == false">
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
              <router-link :to="{ name: 'test-details', params: { id: run.testId } }">{{ tests[run.testId].name }}
              </router-link>
            </td>
            <td class="text-left">
              <router-link :to="{ name: 'run-details', params: { id: run.id } }">view run</router-link>
            </td>
            <td class="text-right">
              <human-date :date="run.startedAt" :relative="false"></human-date>
            </td>
            <td class="text-right">
              <human-date :date="run.finishedAt" :relative="false" v-if="run.startedAt != null && run.finishedAt != null"></human-date>
            </td>
            <td class="text-right">
              <test-result-badge :result="run.result"></test-result-badge>
            </td>
            <td class="text-right">
              <router-link :to="{ name: 'runner-details', params: { id: run.runnerId } }">{{ 'runnerId' in run ?
                  runners[run.runnerId].name : ''
              }}</router-link>
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
      runners: {},
      isLoading: true
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
      this.isLoading = false
    },

  }
}
</script>
