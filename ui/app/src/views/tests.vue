<script setup>
import TestSummary from '../components/TestSummary.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="tests">
            <div style="" v-for="test in Tests" :key="test.id">
                <test-summary :id="test.id" :name="test.name"></test-summary>
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