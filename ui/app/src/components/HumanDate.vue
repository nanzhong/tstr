<template>

    <span>{{ humanDuration }}
        <q-tooltip>
            {{ tooltip }}
        </q-tooltip>
    </span>

</template>


<script>
import * as dayjs from 'dayjs'
import * as duration from 'dayjs/plugin/duration'
import * as relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(duration)
dayjs.extend(relativeTime)


export default {
    props: ['date', 'relative','diff'],
    computed: {
        tooltip() {
            if (this.relative) {
                return this.date.toString()
            }
            return dayjs(this.date).fromNow()
        },
        humanDuration() {
            if (this.date == null) 
                return null;
            if (this.relative) {
                if (this.diff != null) {
                    return dayjs(this.date).from(this.diff)
                }
                return dayjs(this.date).fromNow()
            }
            return this.date.toString()
        }
    }
}
</script>