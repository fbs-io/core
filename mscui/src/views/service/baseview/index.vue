<!--
 * @Author: reel
 * @Date: 2023-06-13 06:16:22
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 05:25:30
 * @Description: 请填写简介
-->
<template>
    <el-main> 
        <el-card v-loading="loading">
            <el-table :data="tableData" style="width: 100%" hight="10rem">
                <el-table-column prop="service" label="服务名称"  >
                    
                </el-table-column>
                <el-table-column prop="node1" label="节点1">

                    <template #default="scope">
                    <div style="color: green;" v-if="scope.row.node1=='1'">
                            <svg width="32" height="32" fill="currentColor" class="bi bi-check-square-fill" viewBox="0 0 16 16">
                                <path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2zm10.03 4.97a.75.75 0 0 1 .011 1.05l-3.992 4.99a.75.75 0 0 1-1.08.02L4.324 8.384a.75.75 0 1 1 1.06-1.06l2.094 2.093 3.473-4.425a.75.75 0 0 1 1.08-.022z"/>
                            </svg>
                        </div>
                        <div style="color: grey;"  v-else>
                            <svg  width="32" height="32" fill="currentColor" class="bi bi-x-square-fill" viewBox="0 0 16 16">
                                <path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2zm3.354 4.646L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 1 1 .708-.708z"/>
                                </svg>
                        </div>
                    </template>    
                </el-table-column>
            </el-table>
         </el-card>
    </el-main>
</template>
<script>
    export default{
        name:"baseview",
        emits:['on-mounted'],
        data(){
            return{
                tableData: [
                    {
                        service: 'ams',
                        node1: '1',
                        
                    },
                    {
                        service: 'cron',
                        node1: '1',
                    },
                    {
                        service: 'msc',
                        node1: '1',
                    },
                    {
                        service: 'rdb',
                        node1: '1',
                    },
                    {
                        service: 'cache',
                        node1: '1',
                    },
                    {
                        service: 'cache',
                        node1: '1',
                    },
                ],
                loading: true
            }
        },
        mounted(){
            this.$emit('on-mounted')
            setInterval(() => {
                if(location.pathname.includes("/service")){
                    this.getSrvStatus()
                    this.loading = false
                }
            }, 5000)
        },
        methods:{
            async  getSrvStatus(){
                var res = await this.$API.overview.srvstatus.get()
                this.tableData= res.details
                
            },
        }
    }

</script>