<script setup lang="ts">
import { Test } from '../api/common/v1/common.pb';

import { defineAsyncComponent } from 'vue'
const HumanDate = defineAsyncComponent(() => import('../components/HumanDate.vue'))

const props = defineProps<{
    test?: Test

}>()

</script>

<template>
    <q-list v-if="test">
        <q-item>
            <q-item-section>
                <q-item-label caption>Name</q-item-label>
                <q-item-label>
                    <router-link :to="{ name: 'test-details', params: { id: test.id } }">{{ test.name }}</router-link>
                </q-item-label>
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
                <ul style="margin-top: 0; margin-bottom: 0;">
                    <li v-for="arg in test.runConfig.args">{{ arg }}</li>
                </ul>
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

        <q-item>
            <q-item-section>
                <q-item-label caption>Registered</q-item-label>
                <q-item-label>
                    <human-date :date="test.registerdAt" :relative="false"></human-date>
                </q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Next run</q-item-label>
                <q-item-label>
                    <human-date :date="test.nextRunAt" :relative="true"></human-date>
                </q-item-label>
            </q-item-section>
        </q-item>

        <q-item>
            <q-item-section>
                <q-item-label caption>Updated</q-item-label>
                <q-item-label>
                    <human-date :date="test.updatedAt" :relative="true"></human-date>
                </q-item-label>
            </q-item-section>
        </q-item>

    </q-list>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent({

})
</script>
