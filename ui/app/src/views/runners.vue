
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
                        <td class="text-right"><span v-if="runner.registeredAt"></span>{{
                                runner.registeredAt
                        }}
                        </td>
                        <td class="text-right"><span v-if="runner.lastHeartbeatAt"></span>{{
                                runner.lastHeartbeatAt
                        }}
                        </td>
                    </tr>
                </tbody>
            </q-markup-table>
        </q-tab-panel>
    </q-page>
</template>


<script>
import tstr from '../tstr'

export default {
    created() {
        this.fetchRunners()
    },
    data() {
        return {
            runners: null,
        }
    },
    methods: {
        async fetchRunners() {
            this.runners = await tstr.fetchRunners()
        }
    }
}
</script>
