<!--
 * @Author: reel
 * @Date: 2023-06-10 19:46:46
 * @LastEditors: reel
 * @LastEditTime: 2023-06-11 15:24:10
 * @Description: 请填写简介
-->
<template>
    <hostinfo></hostinfo>
    <sysinfo></sysinfo>
</template>
<script>
    import hostinfo from "./components/hostinfo.vue"
    import sysinfo from "./components/sysinfo.vue"
    export default {
        emits:['on-mounted'],
        components:{
            hostinfo,
            sysinfo
        },
        mounted(){
            this.$emit('on-mounted')
        },
    }
</script>