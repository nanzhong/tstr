<template>

    <span v-if="date">{{ humanDuration }}
        <q-tooltip>
            {{ tooltip }}
        </q-tooltip>
    </span>

</template>


<script lang="ts">
import { defineComponent } from 'vue';

import dayjs, { Dayjs } from 'dayjs'
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(duration)
dayjs.extend(relativeTime)


export default defineComponent({
    created() {

    },
    props: ['date','relative','diff'],
    // props: {
    //     date: {
    //         type: [dayjs.Dayjs, String],
    //         required: true,
    //     },
    //     relative: Boolean,
    //     diff: {
    //         type: [dayjs.Dayjs, String],
    //     }
    // },
    methods: {
        nDate() {
            return dayjs(this.date)
        },
        nDiff() {
            return  dayjs(this.diff)
        }
    },
    computed: {
        humanDuration(): string {
            if (this.relative) {
                if (this.nDiff()) {
                    return this.nDate().from(this.nDiff())
                }
                return this.nDate().fromNow()
            }
            return this.nDate().toString()
        },

        tooltip(): string {
            if (this.relative) {
                return this.nDate().toString()
            }
            return this.nDate().fromNow()
        }
    }
})

</script>