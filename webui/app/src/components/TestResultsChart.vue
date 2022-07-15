
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
                pass: '#009900',
                fail: '#ff0000',
                error: '#ff9900',
                unknown: '#000000',
            }


            return ['pass', 'fail', 'unknown', 'error'].map(result => {
                return {
                    name: result,
                    color: colors[result],
                    data: this.runs.filter(s => s.Result == result).filter(s => s.StartedAt != null).map(s => {
                        return {
                            x: s.StartedAt.ts,
                            y: Interval.fromDateTimes(s.StartedAt, s.FinishedAt).toDuration(['seconds']).seconds,
                            id: s.ID,
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