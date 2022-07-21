
<template>
    <apexchart type="scatter" height="350" :options="chartOptions" :series="series" @markerClick="openTestRun">
    </apexchart>
</template>

<script>
import tstr from '../tstr'

import { Interval, Duration } from 'luxon'

export default {
    created() {
    },
    data() {
        return {
            chartOptions: {
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
    computed: {
        series() {

            const colors = {
                PASS: '#009900',
                FAIL: '#ff0000',
                ERRROR: '#ff9900',
                UNKNOWN: '#000000',
            }

            return ['PASS', 'FAIL', 'UNKNOWN', 'ERROR'].map(result => {
                return {
                    name: result,
                    color: colors[result],
                    data: this.runs.filter(s => s.result == result).filter(s => s.startedAt != null).map(s => {
                        return {
                            x: s.startedAt.ts,
                            y: Interval.fromDateTimes(s.startedAt, s.finishedAt).toDuration(['seconds']).seconds,
                            id: s.id,
                        }
                    })
                }
            })
        }
    },
    props: ['runs'],
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
    }

};
</script>
