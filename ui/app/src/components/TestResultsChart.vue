
<template>
    <apexchart type="scatter" height="350" :options="chartOptions" :series="series" @markerClick="openTestRun">
    </apexchart>
</template>

<script>
import tstr from '../tstr'

import * as dayjs from 'dayjs'
import * as duration from 'dayjs/plugin/duration'
import * as relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(duration)
dayjs.extend(relativeTime)

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
                yaxis: {
                    decimalsInFloat:1,
                },
                xaxis: {
                    type: 'datetime',
                },
                tooltip: {
                    enabled: true,
                    custom: function ({ series, seriesIndex, dataPointIndex, w }) {
                        const data = w.globals.initialSeries[seriesIndex].data[dataPointIndex];
                        let title = w.globals.tooltip.tooltipTitle.outerHTML;

                        var items = []

                        items.push(`result: ${data.result} <br/>`)
                        if (typeof data.resultData !== 'undefined') {
                            for (const [k, v] of Object.entries(data.resultData)) {
                                items.push(`${k}: ${v} <br/>`)
                            }
                        }

                        return title + '<div>' + items.join("\n") + '</div>'
                    }
                },
                chart: {
                    height: 350,
                    type: 'scatter',
                    zoom: {
                        enabled: true,
                        autoScaleYaxis: true,
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
                            x: s.startedAt,
                            y: dayjs.duration(s.finishedAt.diff(s.startedAt)).asSeconds(),
                            id: s.id,
                            result: s.result,
                            resultData: s.resultData,
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
