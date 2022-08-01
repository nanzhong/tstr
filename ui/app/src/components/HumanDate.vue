<script setup lang="ts">
const props = defineProps<{
    date?: string | object,
    relative: boolean,
    diff?: string | object,
}>()
</script>

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
            if (!this.date) {
                return ""
            }
            if (this.relative) {
                if (this.nDiff()) {
                    return this.nDate().from(this.nDiff())
                }
                return this.nDate().fromNow()
            }
            return this.nDate().toString()
        },

        tooltip(): string {
            if (!this.date) {
                return ""
            }
            if (this.relative) {
                return this.nDate().toString()
            }
            return this.nDate().fromNow()
        }
    }
})

</script>