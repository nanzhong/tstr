<template>
    <q-list v-if="runner != null">
        <q-item>
            <q-item-section>
                <q-item-label caption>Name</q-item-label>
                <q-item-label>{{ runner.Name }}</q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Accept Label Selectors</q-item-label>
                <q-item-label>
                    <q-badge color="gray" v-for="lbl, val in runner.AcceptTestLabelSelectors">
                        {{ lbl }}={{ val }}
                    </q-badge>
                </q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Reject Label Selectors</q-item-label>
                <q-item-label>
                    <q-badge color="gray" v-for="lbl, val in runner.RejectTestLabelSelectors">
                        {{ lbl }}={{ val }}
                    </q-badge>
                </q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Registered</q-item-label>
                <q-item-label>{{ $filters.absoluteDate(runner.RegisteredAt) }}</q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Last Heartbeat</q-item-label>
                <q-item-label>{{ $filters.relativeDate(runner.LastHeartbeatAt) }}</q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
</template>

<script>
export default {
    props: ['runner'],
}
</script>