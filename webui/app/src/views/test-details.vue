<script setup>
import TestDetails from '../components/TestDetails.vue'
</script>

<template>
    <q-page>
        <q-tab-panel name="tests" v-if="test != null">
            <div class="row">
                <div>
                    <div class="text-h6">Test: {{ test.Name }}</div>
                    <test-details :test="test"></test-details>
                </div>
            </div>

            <apexchart type="scatter" height="350" :options="chartOptions" :series="series" @markerClick="openTestRun">
            </apexchart>
        </q-tab-panel>
    </q-page>
</template>


<script>
import tstr from '../tstr'
import { Interval, Duration } from 'luxon'
import { routerKey } from 'vue-router';

export default {
    created() {
        this.fetchTestDetails(this.$route.params.id)
    },
    data() {
        return {
            test: null,
            series: [],
            chartOptions: {
                // tooltip: {
                //     custom: function ({ series, seriesIndex, dataPointIndex, w }) {
                //         const data = w.globals.initialSeries[seriesIndex].data[dataPointIndex];
                //         return `<div class="arrow_box">
                //         <a href="">view test run</a>
                //         </div>`
                //     }
                // },
                dataLabels: {
                    enabled: false
                },
                grid: {
                    xaxis: {
                        lines: {
                            show: true
                        }
                    }
                },
                xaxis: {
                    type: 'datetime',
                },
                chart: {
                    height: 350,
                    type: 'scatter',
                    zoom: {
                    },
                }
            }
        }
    },
    methods: {
        openTestRun(e, ctx, { w, seriesIndex, dataPointIndex }) {
            const data = w.globals.initialSeries[seriesIndex].data[dataPointIndex];
            this.$router.push({
                name: 'run-details',
                params: {
                    id: data.id
                }
            })
        },
        async fetchTestDetails(testId) {
            const testDetails = await tstr.fetchTestDetails(testId)
            this.test = testDetails

            const colors = {
                pass: '#009900',
                fail: '#ff0000',
                error: '#ff9900',
                unknown: '#000000',
            }

            this.series = ['pass', 'fail', 'unknown', 'error'].map(result => {
                return {
                    name: result,
                    color: colors[result],
                    data: testDetails.RunsSummary.filter(s => s.Result == result).filter(s => s.StartedAt != null).map(s => {
                        return {
                            x: s.StartedAt.ts,
                            y: Interval.fromDateTimes(s.StartedAt, s.FinishedAt).toDuration(['seconds']).seconds,
                            id: s.ID,
                        }
                    })
                }
            })
        }
    }
}

</script>