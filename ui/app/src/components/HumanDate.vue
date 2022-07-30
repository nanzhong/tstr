<template>

    <span v-if="date">{{ humanDuration }}
        <q-tooltip>
            {{ tooltip }}
        </q-tooltip>
    </span>

</template>


<script lang="ts">
import { defineComponent } from 'vue';

import dayjs from 'dayjs'
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(duration)
dayjs.extend(relativeTime)


export default defineComponent({
    props: {
        date: {
            type: dayjs.Dayjs,
            required: true,
        },
        relative: Boolean,
        diff: dayjs.Dayjs,
    },
    computed: {
        humanDuration(): string {
            if (this.relative) {
                if (this.diff != null) {
                    return this.date.from(this.diff)
                }
                return this.date.fromNow()
            }
            return this.date.toString()
        },

        tooltip(): string {
            if (this.relative) {
                return this.date.toString()
            }
            return this.date.fromNow()
        }


    }
})

</script>