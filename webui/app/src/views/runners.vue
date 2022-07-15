
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
                <tbody>
                    <tr v-for="runner in this.Runners">
                        <td class="text-left">
                            <router-link :to="{ name: 'runner-details', params: { id: runner.ID } }">{{ runner.Name }}
                            </router-link>
                        </td>
                        <td class="text-right">
                            <q-badge color="gray" v-for="lbl, val in runner.AcceptTestLabelSelectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right">
                            <q-badge color="gray" v-for="lbl, val in runner.RejectTestLabelSelectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right"><span v-if="runner.RegisteredAt.Valid"></span>{{
                                runner.RegisteredAt.Time
                        }}
                        </td>
                        <td class="text-right"><span v-if="runner.LastHeartbeatAt.Valid"></span>{{
                                runner.LastHeartbeatAt.Time
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
            Runners: [],
        }
    },
    methods: {
        async fetchRunners() {
            this.Runners = await tstr.fetchRunners()
        }
    }
}
</script>