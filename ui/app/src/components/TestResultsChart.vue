<script setup lang="ts">
import { RunSummary } from '../api/data/v1/data.pb';

const props = defineProps<{
    runs: RunSummary[]
}>()
</script>

<template>
    <div ref="plot"></div>
</template>

<script lang="ts">
import Plotly from 'plotly.js-dist/plotly'
import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(duration)
dayjs.extend(relativeTime)

export default {
    mounted() {
        const colors = new Map(Object.entries({
            PASS: '#009900',
            FAIL: '#ff0000',
            ERRROR: '#ff9900',
            UNKNOWN: '#000000',
        }))

        const series = ['PASS', 'FAIL', 'UNKNOWN', 'ERROR'].map(result => {

            const points = this.runs.filter(s => s.result == result).filter(s => s.startedAt)

            return {
                name: result.toLowerCase(),
                // hovertemplate: '<b>%{x}</b><br>%{text}<extra></extra>',
                hoverinfo: 'text',
                hovertext: points.map(p => {

                    var items = []

                    items.push(`<b>${p.result}</b><br>`)

                    if (typeof p.resultData !== 'undefined') {
                        for (const [k, v] of Object.entries(p.resultData)) {
                            items.push(`<b>${k}</b>: ${v}`)
                        }
                    }
                    return items.join("<br>")

                }),
                mode: 'markers',
                type: 'scatter',
                marker: {
                    color: colors.get(result),
                },
                result: points.map(p => p.result?.toLowerCase()),
                x: points.map(p => dayjs(p.startedAt).unix() * 1000),
                y: points.map(p => dayjs.duration(dayjs(p.finishedAt).diff(p.startedAt)).asSeconds()),
                id: points.map(p => p.id)
            }
        })


        const layout = {
            xaxis: {
                type: 'date'
            }

        }

        Plotly.newPlot(this.$refs.plot, series, layout)

        this.$refs.plot.on('plotly_click', (data) => {
            const id = data.points[0].data.id[data.points[0].pointIndex]
            this.$router.push({
                name: 'run-details',
                params: {
                    id
                }
            })
        })

    },
};
</script>
