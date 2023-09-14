<!--
 * @Author: reel
 * @Date: 2023-05-30 20:49:47
 * @LastEditors: reel
 * @LastEditTime: 2023-05-30 21:00:33
 * @Description: 请填写简介
-->
<template>
	<el-container>
		<el-header style="height:50px;">
			<div class="common-header-left">
				<div class="common-header-logo">
					<!-- <img :alt="$CONFIG.APP_NAME" src="/img/logo.png"> -->
					<!-- <label>{{$CONFIG.APP_NAME}}</label> -->
				</div>
				<div class="common-header-title">{{title}}</div>
			</div>

		</el-header>
		<el-main>
			<div class="common-container">
				<h2 class="common-title">{{title}}</h2>
				<div class="common-main el-card">
					<slot></slot>
				</div>
			</div>
		</el-main>
	</el-container>
</template>

<script>
	export default {
		props: {
			title: { type: String, default: "" }
		}
	}
</script>

<style>
</style>
