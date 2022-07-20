
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
                            <q-badge color="gray" v-for="lbl, val in runner.accept_test_label_selectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right">
                            <q-badge color="gray" v-for="lbl, val in runner.reject_test_label_selectors">
                                {{ lbl }}={{ val }}
                            </q-badge>
                        </td>
                        <td class="text-right"><span v-if="runner.registered_at"></span>{{
                                runner.registered_at
                        }}
                        </td>
                        <td class="text-right"><span v-if="runner.last_heartbeat_at"></span>{{
                                runner.last_heartbeat_at
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