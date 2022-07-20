<script setup>
import TestSummary from '../components/TestSummary.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="tests">
            <div style="" v-for="test in Tests" :key="test.ID">
                <test-summary :ID="test.ID" :Name="test.Name"></test-summary>
            </div>
        </q-tab-panel>
    </q-page>
</template>

<script>
import tstr from '../tstr'

export default {
    components: {
        TestSummary,
    },
    created() {
        this.fetchTests()
    },
    data() {
        return {
            Tests: [],
        };
    },
    mounted() {
    },
    methods: {
        async fetchTests() {
            this.Tests = await tstr.fetchTests()
        }
    }
}

</script>

<style>
</style>