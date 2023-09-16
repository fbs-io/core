<!--
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2023-06-14 21:33:07
 * @Description: 请填写简介
-->
<template>
	<el-main>
		<el-card shadow="hover" header="基础信息" >
			<div >
				<el-row style="font-size: 14px;">
					<el-col :span="4" :size="50">
						<span >主机名: {{ hostinfoData.host }}</span>
					</el-col>
					<el-col :span="8">
						<span>CPU: {{ hostinfoData.cpu }}</span>
					</el-col>
					<el-col :span="4">
						<span>内存: {{ hostinfoData.mem }} MB </span>
					</el-col>
					<el-col :span="4">
						<span>协程数: {{ hostinfoData.numgoroutine }} 个</span>
					</el-col>
					<el-col :span="4">
						<span>系统: {{ hostinfoData.os }}</span>
					</el-col>

				</el-row>
				<br>
				<el-row style="font-size: 14px;">
					<el-col :span="4" :size="50">
						<span >Golang版本: {{ hostinfoData.goversion }}</span>
					</el-col>
					<el-col :span="4">
						<span>Core版本: {{ hostinfoData.version }}</span>

					</el-col>
					<el-col :span="4">
						<span>APP版本: {{ hostinfoData.appversion }}</span>
					</el-col>
					<el-col :span="6">
						<span>系统运行: {{ hostinfoData.osrunday.split("-")[0] }} 
							<span style="font-size: 10px;"> {{hostinfoData.osrunday.split("-")[1]}}</span>
						</span>
					</el-col>
					<el-col :span="6">
						<span>系统运行: {{ hostinfoData.apprunday.split("-")[0] }} 
							<span style="font-size: 10px;"> {{hostinfoData.apprunday.split("-")[1]}}</span>
						</span>
					</el-col>
				</el-row>
				<br>
				<el-row style="font-size: 14px;">
					<el-col :span="4" :size="50">
						<span >GC次数: {{ hostinfoData.numgc }} 次</span>
					</el-col>
					<el-col :span="4">
						<span>下次GC内存: {{ hostinfoData.nextgc }} MB</span>

					</el-col>
					<el-col :span="4">
						<span>GC总时间: {{ hostinfoData.pausetotal }} 秒</span>

					</el-col>
					<el-col :span="6">
						<span>上次GC耗时: {{ hostinfoData.pause }} 毫秒</span>
					</el-col>
					<el-col :span="6">
						<span>上次GC时间: {{ hostinfoData.lastgc }} 秒</span>
					</el-col>
				</el-row>

			</div>
		</el-card>
	</el-main>
</template>

<script>
	export default {
		name:"baseinfo",
		data() {
			return {
				hostinfoData:{
					host:"",
					cpu:"",
					os:"",
					mem:"",
					numgoroutine:"",
					osrunday:"",
					numgc:"",
					nextgc:"",
					pausetotal:"",
					pause:"",
					apprunday:"",
					appversion:"",
					version:"",
					goversion:"",
				},
			}
		},
		 mounted(){
			this.hostinfo()
		},
		methods: {
			async hostinfo (){
				let  res =  await this.$API.overview.hostinfo.get()
				this.hostinfoData = res.details

			}
		}
	}
</script>

<style scoped>

</style>
