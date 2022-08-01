<script setup lang="ts">
import HumanDate from '../components/HumanDate.vue'
</script>

<template>
    <q-page>

        <q-tab-panel name="runners">
            <div class="text-h6">Runners</div>
            <q-markup-table separator="horizontal" flat bordered>
                <thead>
                    <tr>
                        <th class="text-left">Name</th>
                        <th class="text-right">Accept Selectors</th>
                        <th class="text-right">Reject Selectors</th>
                        <th class="text-right">Registered</th>
                        <th class="text-right">Last Heartbeat</th>
                    </tr>
                </thead>
                <tbody v-if="runners != null">
                    <tr v-for="runner in runners">
                        <td class="text-left">
                            <router-link :to="{ name: 'runner-details', params: { id: runner.id } }">{{ runner.name }}
                            </router-link>
                        </td>
                        <td class="text-right">
                            <q-badge color="gray" v-for="lbl, val in runner.acceptTestLabelSelectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right">
                            <q-badge color="gray" v-for="lbl, val in runner.rejectTestLabelSelectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right">
                            <human-date :date="runner.registeredAt" :relative="false"></human-date>
                        </td>
                        <td class="text-right">
                            <human-date :date="runner.lastHeartbeatAt" :relative="true"></human-date>
                        </td>
                    </tr>
                </tbody>
            </q-markup-table>
        </q-tab-panel>
    </q-page>
</template>


<script lang="ts">
import { Runner } from '../api/common/v1/common.pb';
import { DataService } from '../api/data/v1/data.pb';
import { InitReq } from '../api/fetch.pb';

const initReq: InitReq = {
    pathPrefix: '/api'
}

export default {
    created() {
        (async () => {
            const runnersResponse = await DataService.QueryRunners({}, initReq);
            if (runnersResponse.runners) {
                this.runners = runnersResponse.runners
            }

        })();
    },
    data() {
        return {
            runners: [] as Runner[],
        }
    },
}
</script>
