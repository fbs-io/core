<!--
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2023-06-12 06:25:44
 * @Description: 请填写简介
-->
<template>
	<el-table :data="tableData" style="width: 100%">
	  <el-table-column label="Date" width="180">
		<template #default="scope">
		  <div style="display: flex; align-items: center">
			<el-icon><timer /></el-icon>
			<span style="margin-left: 10px">{{ scope.row.date }}</span>
		  </div>
		</template>
	  </el-table-column>
	  <el-table-column label="Name" width="180">
		<template #default="scope">
		  <el-popover effect="light" trigger="hover" placement="top" width="auto">
			<template #default>
			  <div>name: {{ scope.row.name }}</div>
			  <div>address: {{ scope.row.address }}</div>
			</template>
			<template #reference>
			  <el-tag>{{ scope.row.name }}</el-tag>
			</template>
		  </el-popover>
		</template>
	  </el-table-column>
	  <el-table-column label="Operations">
		<template #default="scope">
		  <el-button size="small" @click="handleEdit(scope.$index, scope.row)"
			>Edit</el-button
		  >
		  <el-button
			size="small"
			type="danger"
			@click="handleDelete(scope.$index, scope.row)"
			>Delete</el-button
		  >
		</template>
	  </el-table-column>
	</el-table>
  </template>
  
  <script setup>
  import { Timer } from '@element-plus/icons-vue'
  
  
  const handleEdit = (index, row) => {
	console.log(index, row)
  }
  const handleDelete = (index, row) => {
	console.log(index, row)
  }
  
  const tableData = [
	{
	  date: '2016-05-03',
	  name: 'Tom',
	  address: 'No. 189, Grove St, Los Angeles',
	},
	{
	  date: '2016-05-02',
	  name: 'Tom',
	  address: 'No. 189, Grove St, Los Angeles',
	},
	{
	  date: '2016-05-04',
	  name: 'Tom',
	  address: 'No. 189, Grove St, Los Angeles',
	},
	{
	  date: '2016-05-01',
	  name: 'Tom',
	  address: 'No. 189, Grove St, Los Angeles',
	},
  ]
  </script>
  

<style scoped>
	.item-background p {color: #999;margin-top:10px;line-height: 1.8;}
</style>
