<!--
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2024-03-27 05:36:57
 * @Description: 请填写简介
-->
<template>
	<el-main>
	<el-row :gutter="24">
    	<el-col :span="14">
			<el-card shadow="hover" header="系统状态"  style="height: 20rem;" >
				<div class="progress" >
					<el-progress type="dashboard" :percentage="Number(info.cpu)" :width="progressWidth" :color="colors">
						<template #default="{ percentage }">
							<div class="percentage-value">{{ percentage }}%</div>
							<div class="percentage-label">CPU使用率</div>
						</template>
					</el-progress>
					
					<el-progress type="dashboard" :percentage="Number(info.memp)" :width="progressWidth" :color="colors">
						<template #default="{ percentage }">
							<div class="percentage-value">{{ percentage }}%</div>
							<div class="percentage-label">内存使用率</div>
						</template>
					</el-progress>
					<el-progress type="dashboard" :percentage="Number(info.appmemp)" :width="progressWidth"  :color="colors">
						<template #default="{ percentage }">
							<div class="percentage-value">{{ percentage }}%</div>
							<div class="percentage-label">APP内存使用率</div>
						</template>
					</el-progress>
					<el-progress type="dashboard" :percentage="Number(info.disk)" :width="progressWidth"  :color="colors">
						<template #default="{ percentage }">
							<div class="percentage-value">{{ percentage }}%</div>
							<div class="percentage-label">磁盘使用率</div>
						</template>
					</el-progress>
				</div>
			</el-card>
		</el-col>
		<el-col :span="10">
			<el-card shadow="hover" header="CPU(%)" v-loading="loading" style="height: 20rem;">
				<scEcharts  ref="c1" height="12rem" :option="cpuOption"></scEcharts>
			</el-card>
		</el-col>
	</el-row>
	<el-row :gutter="24">
		<el-col :span="14">
			<el-card shadow="hover" header="进程信息" v-loading="loading" style="height: 41.25rem;">
				
				<el-table
					:data="processInfo"
					style="width: 100%"
					height = "40rem"
				>
					<el-table-column prop="pid" label="PID" width="80" />
					<el-table-column prop="pname" label="进程名称" width="260" />
					<el-table-column prop="cpupercent" label="CPU" />
					<el-table-column prop="meminfo" label="内存" />
					<el-table-column prop="io" label="IO" />
				</el-table>
			</el-card>
		</el-col>
		<el-col :span="10">
			<el-card shadow="hover" header="内存(%)" v-loading="loading" style="height: 20rem;">
				<scEcharts  ref="c1" height="12rem" :option="memOption"></scEcharts>
			</el-card>
			<el-card shadow="hover" header="APP(MB)" v-loading="loading" style="height: 20rem;">
				<scEcharts  ref="c1" height="12rem" :option="appMemOption"></scEcharts>
			</el-card>
		</el-col>
	</el-row>
</el-main>
</template>

<script>
	import scEcharts from '@/components/scEcharts';
	export default {
		
		data() {
			return {
				progressWidth:120,
				colors : [
					{ color: '#5cb87a', percentage: 25 },
					{ color: '#6f7ad3', percentage: 60 },
					{ color: '#e6a23c', percentage: 85 },
					{ color: '#f56c6c', percentage: 100 },
				],
				info:{
					cpu:0,
					memp:0,
					appmemp:0,
					disk:0,
				},
				cpuOption:{
					tooltip: {
						trigger: 'axis'
					},
					grid:{
						x:5,
						y:5,
						x2:5,
						y2:2
					},

					xAxis: {
						boundaryGap: false,
						type: 'category',
						data:[],
					},
					yAxis: [{
						type: 'value',
						// name: '百分比',
						"splitLine": {
							"show": true
						}
					}],
					series: [
						{
							name: 'CPU 使用率',
							type: 'line',
							symbol: 'none',
							lineStyle: {
								width: 1,
								color: '#409EFF'
							},
							areaStyle: {
								opacity: 0.1,
								color: '#79bbff'
							},
							data:[],
						},
					],
				},
				memOption:{
					tooltip: {
						trigger: 'axis'
					},
					grid:{
						x:5,
						y:5,
						x2:5,
						y2:5
					},

					xAxis: {
						boundaryGap: false,
						type: 'category',
						data:[],
					},
					yAxis: [{
						type: 'value',
						// name: '百分比',
						"splitLine": {
							"show": true
						}
					}],
					series: [
						{
							name: '内存使用率',
							type: 'line',
							symbol: 'none',
							lineStyle: {
								width: 1,
								color: '#409EFF'
							},
							areaStyle: {
								opacity: 0.1,
								color: '#79bbff'
							},
							data:[],
						},
					],
				},
				appMemOption:{
					tooltip: {
						trigger: 'axis'
					},
					grid:{
						x:5,
						y:5,
						x2:5,
						y2:5
					},

					xAxis: {
						boundaryGap: false,
						type: 'category',
						data:[],
					},
					yAxis: [{
						type: 'value',
						// name: '百分比',
						"splitLine": {
							"show": true
						}
					}],
					series: [
						{
							name: 'APP内存占用',
							type: 'line',
							symbol: 'none',
							lineStyle: {
								width: 1,
								color: '#409EFF'
							},
							areaStyle: {
								opacity: 0.1,
								color: '#79bbff'
							},
							data:[],
						},
					],
				},
				loading:true,
				processInfo:[],

			}
		},
		mounted(){
				setInterval(() => {
					if(location.pathname.includes("/overview")){
						this.sys()
						this.loading = false
					}
				}, 5000)			
		},
		methods:{
			async sys(){
				var res = await this.$API.overview.sysinfo.get()
				this.info = res.details
				this.cpuOption.series[0].data.push(res.details.cpu)
				this.memOption.series[0].data.push(res.details.memp)
				this.appMemOption.series[0].data.push(res.details.appmem)
				if (this.cpuOption.series[0].data.length>120){
					this.cpuOption.series[0].data.shift()
					this.memOption.series[0].data.shift()
					this.appMemOption.series[0].data.shift()
				}
				var now = new Date();
				this.cpuOption.xAxis.data.push(now.toLocaleTimeString().replace(/^\D*/,''))
				this.memOption.xAxis.data.push(now.toLocaleTimeString().replace(/^\D*/,''))
				this.appMemOption.xAxis.data.push(now.toLocaleTimeString().replace(/^\D*/,''))
				if (this.cpuOption.xAxis.data.length>120){
					this.cpuOption.xAxis.data.shift()
					this.memOption.xAxis.data.shift() 
					this.appMemOption.xAxis.data.shift() 
				}
				this.processInfo = res.details.processinfo
			},

			// async processinfo(){
			// 	var res = await this.$API.overview.processinfo.get(){}
			// }
		}
	}
</script>

<style scoped>
	.progress {text-align: center; }
	.progress .el-progress {
  		margin-right: 80px;
	}
	.progress .percentage-value {font-size: 28px;}
	.progress .percentage-label {font-size: 12px;margin-top: 10px;}
</style>
