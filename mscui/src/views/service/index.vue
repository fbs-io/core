<!--
 * @Author: reel
 * @Date: 2023-06-12 07:23:49
 * @LastEditors: reel
 * @LastEditTime: 2023-06-13 07:59:09
 * @Description: 请填写简介
-->
<template>
	<div v-if="pageLoading">
		<el-main>
			<el-card shadow="never">
				<el-skeleton :rows="1"></el-skeleton>
			</el-card>
			<el-card shadow="never" style="margin-top: 15px;">
				<el-skeleton></el-skeleton>
			</el-card>
		</el-main>
	</div>
    
	<baseview @on-mounted="onMounted"></baseview>
</template>

<script>
	import { defineAsyncComponent } from 'vue';
	const baseview = defineAsyncComponent(() => import('./baseview'));
	export default {
		name: "overview",
		components: {
			baseview
		},
		data(){
			return {
				pageLoading: true,
			}
		},
		created(){
		},
		mounted(){
		},
		methods: {
			onMounted(){
				this.pageLoading = false
			}
		}
	}
</script>

<style>
</style>
