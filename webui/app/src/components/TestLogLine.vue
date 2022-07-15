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
            return atob(this.logline.data)
        },
        style() {
            var borderColor = ""

            if (this.logline.type == "STDOUT")
                borderColor = "red"
            else
                borderColor = "green"
            return `padding-left:0.5em; border-left: 2px solid ${borderColor}`

        }

    },
}
</script>