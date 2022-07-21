<template>
    <q-list v-if="runner != null">
        <q-item>
            <q-item-section>
                <q-item-label caption>Name</q-item-label>
                <q-item-label>{{ runner.name }}</q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Accept Label Selectors</q-item-label>
                <q-item-label>
                    <q-badge color="gray" v-for="lbl, val in runner.acceptTestLabelSelectors">
                        {{ lbl }}={{ val }}
                    </q-badge>
                </q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Reject Label Selectors</q-item-label>
                <q-item-label>
                    <q-badge color="gray" v-for="lbl, val in runner.rejectTestLabelSelectors">
                        {{ lbl }}={{ val }}
                    </q-badge>
                </q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Registered</q-item-label>
                <q-item-label>{{ $filters.absoluteDate(runner.registeredAt) }}</q-item-label>
            </q-item-section>
        </q-item>
        <q-item>
            <q-item-section>
                <q-item-label caption>Last Heartbeat</q-item-label>
                <q-item-label>{{ $filters.relativeDate(runner.lastHeartbeatAt) }}</q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
</template>

<script>
export default {
    props: ['runner'],
}
</script>
