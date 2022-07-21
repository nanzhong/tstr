<script setup>
import moment from 'moment'
</script>

<template>
    <q-list v-if="test != null">
        <q-item>
            <q-item-section>
                <q-item-label caption>Name</q-item-label>
                <q-item-label>{{ test.name }}</q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Container Image</q-item-label>
                <q-item-label>
                    {{ test.runConfig.containerImage }}
                </q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Command</q-item-label>
                <q-item-label>
                    {{ test.runConfig.command }}
                </q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Args</q-item-label>
                <q-item-label>
                    <ul>
                        <li v-for="arg in test.runConfig.args">{{ arg }}</li>
                    </ul>
                </q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Cron Schedule</q-item-label>
                <q-item-label>{{ test.cronSchedule }}</q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Labels</q-item-label>
                <q-item-label>
                    <q-badge color="gray" v-for="lbl, val in test.labels">
                        {{ lbl }}={{ val }}
                    </q-badge>
                </q-item-label>
            </q-item-section>
        </q-item>


        <!-- <q-item>
            <q-item-section>
                <q-item-label caption>Created</q-item-label>
                <q-item-label>{{ $filters.absoluteDate(test.createdAt) }}</q-item-label>
            </q-item-section>
        </q-item> -->

        <q-item>
            <q-item-section>
                <q-item-label caption>Registered</q-item-label>
                <q-item-label>{{ $filters.absoluteDate(test.registeredAt) }}</q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Next run</q-item-label>
                <q-item-label>{{ $filters.relativeDate(test.nextRunAt) }}</q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Updated</q-item-label>
                <q-item-label>{{ $filters.absoluteDate(test.updatedAt) }}</q-item-label>
            </q-item-section>
        </q-item>

    </q-list>
</template>

<script>
export default {
    props: ['test'],
}
</script>
