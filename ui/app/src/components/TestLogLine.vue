<template>
    <span :style="style">
        {{ logline.time }}
        {{ txtData }}
        <br />
    </span>
</template>

<script>
export default {
    props: ['logline'],
    computed: {
        txtData() {
            if (this.logline.data && this.logline.data.length > 0) {
                return atob(this.logline.data)
            }
            return ""
        },
        style() {
            var borderColor = ""

            if (this.logline.type == "STDOUT")
                borderColor = "red"
            else
                borderColor = "green"
            return `font-family: monospace; padding-left:0.5em; border-left: 3px solid ${borderColor}`

        }

    },
}
</script>